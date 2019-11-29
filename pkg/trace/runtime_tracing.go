package trace

import (
	"context"
	t "runtime/trace"
)

func NewTask(name string) func() {
	if !t.IsEnabled() {
		// Avoid additional overhead if
		// runtime/trace is not enabled.
		return func() {}
	}
	region := t.StartRegion(context.Background(), name)
	return region.End
}
