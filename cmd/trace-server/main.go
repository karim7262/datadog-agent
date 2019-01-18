package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DataDog/datadog-agent/pkg/trace/pb"
	"github.com/DataDog/datadog-agent/pkg/trace/test"

	"github.com/tinylib/msgp/msgp"
)

func main() {
	r := test.Runner{Verbose: true}
	if err := r.Start(); err != nil {
		log.Fatal(err)
	}
	defer r.Shutdown(2 * time.Second)

	if err := r.RunAgent([]byte("")); err != nil {
		log.Fatal(err)
	}
	defer r.KillAgent()

	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var traces pb.Traces
		if err := msgp.Decode(req.Body, &traces); err != nil {
			fmt.Println(err)
			return
		}
		if err := r.Post(traces); err != nil {
			fmt.Println(err)
			return
		}
		timeout := time.After(2 * time.Second)
	outer:
		for {
			select {
			case p := <-r.Out():
				if v, ok := p.(pb.TracePayload); ok {
					if err := json.NewEncoder(os.Stdout).Encode(v); err != nil {
						fmt.Println(err)
						return
					}
					break outer
				}
			case <-timeout:
				fmt.Println("timed out")
				break outer
			}
		}
	}))
}
