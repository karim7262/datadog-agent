//+build linux_bpf

package ebpf

import (
	"fmt"
)

// bpfMapName stores the name of the BPF maps storing statistics and other info
type bpfMapName string

const (
	connMap            bpfMapName = "conn_stats"
	tcpStatsMap        bpfMapName = "tcp_stats"
	tcpCloseEventMap   bpfMapName = "tcp_close_events"
	latestTimestampMap bpfMapName = "latest_ts"
	tracerStatusMap    bpfMapName = "tracer_status"
	portBindingsMap    bpfMapName = "port_bindings"
	telemetryMap       bpfMapName = "telemetry"
)

// sectionName returns the sectionName for the given BPF map
func (b bpfMapName) sectionName() string {
	return fmt.Sprintf("maps/%s", b)
}
