//
// how to test:
// - compile: go build -tags "cpython" -gcflags='-N -l' test-race.go
// - run: GOTRACEBACK=crash ./test-race
//   you can also reduce the number of goroutine by disabling Go GC (with
//   GOGC=off), it will still segfault but the cordump has fewer goroutine

// then "gdb ./test-race core" to debug ("info goroutine" and "goroutine #id bt" to debug)

package main

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/pkg/collector/py"
	//python "github.com/sbinet/go-python"
	python "github.com/DataDog/go-python3"
)

func main() {
	aggregator.InitAggregatorWithFlushInterval(nil, "", "", time.Hour)
	state := py.Initialize(
		"/home/hush-hush/dev/datadog-agent/pkg/collector/py/tests/",
		"/home/hush-hush/dev/datadog-agent/bin/agent/dist/",
	)

	l, err := py.NewPythonCheckLoader()
	fmt.Printf("Error: NewPythonCheckLoader: %v\n", err)
	for i := 0; i < 5000; i += 1 {
		fmt.Printf("iter %d \n", i)

		config := integration.Config{Name: "testcheck"}
		config.Instances = append(config.Instances, []byte("foo: bar"))
		config.Instances = append(config.Instances, []byte("bar: baz"))
		module, err := l.Load(config)
		fmt.Printf("load testcheck: %v %v\n", module, err)

		config = integration.Config{Name: "foo"}
		module, err = l.Load(config)
		fmt.Printf("load foo: %v %v\n", module, err)

		fmt.Printf("====================================\n")
	}
	python.PyEval_RestoreThread(state)
}
