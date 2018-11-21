package configtryer

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func status2XXChecker(r *http.Response) bool {
	return r.StatusCode/100 == 2
}

func TestEmptyRules(t *testing.T) {
	name := "redis"

	tryer := RulesTryer{
		map[string]rules{
			name: rules{},
		},
	}

	conf, err := tryer.Try(name)
	assert.Error(t, err)
	assert.Nil(t, conf)
}

func TestValidSearchHTTPPort(t *testing.T) {
	name := "elasticsearch"
	port := 9200
	uri := "_cluster/health"
	var wg sync.WaitGroup

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status": "yellow"}`)
	})

	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("/%s", uri), handler)
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}

	go func() {
		wg.Add(1)
		server.ListenAndServe()
		wg.Done()
	}()

	tryer := RulesTryer{
		map[string]rules{
			name: rules{
				httpEndpoints: []endpoint{{
					uri,
					port,
				}},
				httpChecker: status2XXChecker,
			},
		},
	}

	conf, err := tryer.Try(name)
	server.Shutdown(context.Background())

	assert.NoError(t, err)

	assert.Equal(t, 1, len(conf.Ports))
	assert.Equal(t, 0, len(conf.UnixSockets))
	assert.Equal(t, port, conf.Ports[0])
	wg.Wait()
}

func TestInvalidSearchHTTPPort(t *testing.T) {
	name := "elasticsearch"
	port := 9200
	uri := "_cluster/health"

	tryer := RulesTryer{
		map[string]rules{
			name: rules{
				httpEndpoints: []endpoint{{
					uri,
					port,
				}},
				httpChecker: status2XXChecker,
			},
		},
	}

	conf, err := tryer.Try(name)
	assert.Error(t, err)
	assert.Nil(t, conf)
}

func TestValidSearchSockets(t *testing.T) {
	name := "redis"

	tmpfile, err := ioutil.TempFile("", "redis.sock")
	assert.NoError(t, err)
	file := tmpfile.Name()
	defer os.Remove(file)

	tryer := RulesTryer{
		map[string]rules{
			name: rules{
				socketGlobs: []string{
					"/tmp/redis.socket",
					"/tmp/redis*.sock",
					"/tmp/redis.sock",
					"/tmp/*redis.sock",
					file,
					fmt.Sprintf("%s*", file[:len(file)-5]), // Fake glob
				},
			},
		},
	}

	conf, err := tryer.Try(name)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(conf.Ports))
	// We should not have duplicates
	assert.Equal(t, 1, len(conf.UnixSockets))
	assert.Equal(t, file, conf.UnixSockets[0])
}

func TestInvalidSearchSockets(t *testing.T) {
	name := "redis"

	tryer := RulesTryer{
		map[string]rules{
			name: rules{
				socketGlobs: []string{
					"/tmp/redis.socket",
					"/tmp/redis*.sock",
					"/tmp/redis.sock",
					"/tmp/*redis.sock",
				},
			},
		},
	}

	conf, err := tryer.Try(name)
	assert.Error(t, err)
	assert.Nil(t, conf)
}

func TestIntRange(t *testing.T) {
	r := intRange(5, 10)

	assert.Contains(t, r, 5)
	assert.Contains(t, r, 6)
	assert.Contains(t, r, 7)
	assert.Contains(t, r, 8)
	assert.Contains(t, r, 9)
	assert.Equal(t, 5, len(r))
}
