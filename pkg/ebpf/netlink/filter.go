// +build linux

package netlink

import (
	"bytes"
	"fmt"

	"github.com/DataDog/datadog-agent/pkg/util/log"
	ct "github.com/florianl/go-conntrack"
	bpflib "github.com/iovisor/gobpf/elf"
)

type natFilter struct {
	module *bpflib.Module
	filter *bpflib.SocketFilter
	fds    []int
}

func loadNATFilter() (*natFilter, error) {
	buf, err := Asset("nat_filter-debug.o")
	if err != nil {
		return nil, fmt.Errorf("couldn't find asset: %s", err)
	}

	m := bpflib.NewModuleFromReader(bytes.NewReader(buf))
	if m == nil {
		return nil, fmt.Errorf("BPF not supported")
	}

	err = m.Load(map[string]bpflib.SectionParams{"socket/dns_filter": {}})
	if err != nil {
		return nil, fmt.Errorf("could not load bpf module: %s", err)
	}

	filter := m.SocketFilter("socket/nat_filter")
	if filter == nil {
		return nil, fmt.Errorf("error retrieving socket filter")
	}

	return &natFilter{module: m, filter: filter}, nil
}

func (n *natFilter) Attach(nfct *ct.Nfct) {
	if n == nil {
		return
	}

	rawConn, err := nfct.Con.SyscallConn()
	if err != nil {
		log.Warn("error obtaining syscall.RawConn", err)
		return
	}

	ctrlErr := rawConn.Control(func(socketFD uintptr) {
		if err := bpflib.AttachSocketFilter(n.filter, int(socketFD)); err != nil {
			log.Warn("error attaching BPF filter to socket: %s", err)
			return
		}

		// Keep track of socket FD so we can later detach the filter
		n.fds = append(n.fds, int(socketFD))
	})

	if ctrlErr != nil {
		log.Warn("error executing operation over NETLINK socket: %s", ctrlErr)
	}
}

func (n *natFilter) Close() {
	if n == nil {
		return
	}

	for socketFD := range n.fds {
		if err := bpflib.DetachSocketFilter(n.filter, socketFD); err != nil {
			log.Errorf("error detaching socket filter: %s", err)
		}
	}

	_ = n.module.Close()
}
