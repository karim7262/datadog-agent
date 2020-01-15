// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

package testsuite

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DataDog/datadog-agent/pkg/trace/pb"
	"github.com/DataDog/datadog-agent/pkg/trace/test"
	"github.com/DataDog/datadog-agent/pkg/trace/test/testutil"
	"github.com/stretchr/testify/assert"
)

type receivedPayload struct {
	path    string
	apiKey  string
	payload pb.TracePayload
}

func newBackendEndpoints(t *testing.T) (server *httptest.Server, out <-chan receivedPayload) {
	outCh := make(chan receivedPayload)
	reply := func(endpoint string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			p, err := test.TracePayloadFromRequest(r)
			if err != nil {
				t.Fatal(err)
			}
			outCh <- receivedPayload{
				path:    endpoint,
				payload: p,
				apiKey:  r.Header.Get("Dd-Api-Key"),
			}
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/one/", reply("/one"))
	mux.HandleFunc("/two/", reply("/two"))
	return httptest.NewServer(mux), outCh
}

// TestMultipleEndpoints ensures that the additional_endpoints settings is sending the correct payload to
// all destinations.
func TestMultipleEndpoints(t *testing.T) {
	r := test.Runner{Verbose: true}
	if err := r.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := r.Shutdown(time.Second); err != nil {
			t.Log("shutdown: ", err)
		}
	}()
	srv, out := newBackendEndpoints(t)
	defer srv.Close()
	conf := fmt.Sprintf(`
apm_config:
  additional_endpoints:
    %[1]s/one:
      - apikey11
      - apikey12
    %[1]s/two:
      - apikey2
`, srv.URL)
	if err := r.RunAgent([]byte(conf)); err != nil {
		t.Fatal(err)
	}
	defer r.KillAgent()

	if err := r.Post(pb.Traces{
		testutil.RandomTrace(1, 4),
		testutil.RandomTrace(4, 10),
		testutil.RandomTrace(2, 2),
	}); err != nil {
		t.Fatal(err)
	}
	waitForTrace(t, &r, func(v pb.TracePayload) {
		if n := len(v.Traces); n != 3 {
			t.Fatalf("expected %d traces, got %d", 1, n)
		}
		expect := map[string]bool{
			"/one:apikey11": true,
			"/one:apikey12": true,
			"/two:apikey2":  true,
		}
		var got int
		timeout := time.After(5 * time.Second)
	loop:
		for {
			select {
			case <-timeout:
				t.Fatal("time out")
			case p := <-out:
				got++
				delete(expect, fmt.Sprintf("%s:%s", p.path, p.apiKey))
				assert.Equal(t, p.payload, v)
				if got == 3 {
					if len(expect) > 0 {
						t.Fatalf("didn't hit endpoint(s) %+v", expect)
					}
					break loop
				}
			}
		}
	})
}
