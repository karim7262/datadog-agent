package configtryer

import (
	"net"
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
			tcpChecker: func(conn *net.TCPConn) bool {
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
	},
}
