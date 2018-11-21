package configtryer

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"time"
)

type endpoint struct {
	uri  string
	port int
}

type rules struct {
	tcpPorts      []int
	httpEndpoints []endpoint
	socketGlobs   []string

	tcpChecker  func(*net.TCPConn) bool
	httpChecker func(*http.Response) bool
}

func (r *rules) searchHTTPPorts() []int {
	ports := []int{}

	if r.httpEndpoints != nil && r.httpChecker != nil {
		client := http.Client{Timeout: time.Duration(2 * time.Second)}

		for _, e := range r.httpEndpoints {
			resp, err := client.Get(fmt.Sprintf("http://localhost:%d/%s", e.port, e.uri))
			if err != nil {
				// Probably nothing listening there
				continue
			}

			// If the httpChecker returns true store the port
			if r.httpChecker(resp) {
				ports = append(ports, e.port)
			}
		}
	}

	return ports
}

func (r *rules) searchTCPPorts() []int {
	ports := []int{}

	if r.tcpPorts != nil && r.tcpChecker != nil {
		for _, p := range r.tcpPorts {
			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", p))
			if err != nil {
				// TODO logging, this should not fail
				continue
			}

			conn, err := net.DialTCP("tcp", nil, addr)
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

			if err != nil {
				// Probably nothing listening there
				continue
			}

			// If the tcpChecker returns true store the port
			if r.tcpChecker(conn) {
				ports = append(ports, p)
			}

			conn.Close()
		}
	}

	return ports
}

func (r *rules) searchSockets() []string {
	sockets := []string{}
	seen := map[string]struct{}{}

	if r.socketGlobs != nil {
		for _, g := range r.socketGlobs {
			matches, err := filepath.Glob(g)
			// The only possible returned error for Glob is ErrBadPattern, so there should not be any error
			if err != nil {
				continue
			}

			if matches != nil {
				for _, m := range matches {
					// Assert we don't have duplicates
					if _, ok := seen[m]; !ok {
						sockets = append(sockets, m)
						seen[m] = struct{}{}
					}
				}
			}
		}
	}

	return sockets
}

type RulesTryer struct {
	rulesMap map[string]rules
}

func (rt *RulesTryer) Try(name string) (*Config, error) {
	r, ok := rt.rulesMap[name]

	if !ok {
		return nil, fmt.Errorf("no rules found for %s", name)
	}

	conf := &Config{
		Ports: append(
			r.searchHTTPPorts(),
			r.searchTCPPorts()...,
		),
		UnixSockets: r.searchSockets(),
	}

	if len(conf.Ports) == 0 && len(conf.UnixSockets) == 0 {
		return nil, fmt.Errorf("no ports or sockets found matching the given rules for %s", name)
	}

	return conf, nil
}

func status2XXChecker(r *http.Response) bool {
	return r.StatusCode/100 == 2
}
