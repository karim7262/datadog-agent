package configtryer

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func Try(name string) (*Config, error) {
	return DefaultTryer.Try(name)
}

var DefaultTryer = RulesTryer{
	rulesMap: map[string]rules{
		// Redis rules
		"redisdb": rules{
			socketGlobs: []string{
				"/tmp/redis*sock*",
			},
			tcpPorts: intRange(6375, 6386),
			tcpChecker: func(conn net.Conn) bool {
				if _, err := conn.Write([]byte("echo hi\n")); err != nil {
					return false
				}

				raw := make([]byte, 1024)

				if _, err := conn.Read(raw); err != nil {
					return false
				}

				s := string(raw)
				return strings.HasPrefix(s, "$2") || strings.HasPrefix(s, "-NOAUTH Authentication required.")
			},
		},

		// Elasticsearch rules
		"elastic": rules{
			httpEndpoints: endpointRange("_cluster/health", 9195, 9206),
			httpChecker: func(r *http.Response) bool {
				res := struct{ status string }{}
				defer r.Body.Close()
				raw, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return false
				}

				err = json.Unmarshal(raw, &res)
				return err == nil && res.status != ""
			},
		},
	},
}
