// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package listeners

import (
	"expvar"
	"fmt"
	"net"
	"time"

	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/tidwall/evio"
)

var (
	udpExpvars             = expvar.NewMap("dogstatsd-udp")
	udpPacketReadingErrors = expvar.Int{}
	udpPackets             = expvar.Int{}
)

func init() {
	udpExpvars.Set("PacketReadingErrors", &udpPacketReadingErrors)
	udpExpvars.Set("Packets", &udpPackets)
}

// UDPListener implements the StatsdListener interface for UDP protocol.
// It listens to a given UDP address and sends back packets ready to be
// processed.
// Origin detection is not implemented for UDP.
type UDPListener struct {
	packetPool   *PacketPool
	packetBuffer *packetBuffer
	events       evio.Events
	url          string
	stopped      bool
}

// NewUDPListener returns an idle UDP Statsd listener
func NewUDPListener(packetOut chan Packets, packetPool *PacketPool) (*UDPListener, error) {
	var address string

	if config.Datadog.GetBool("dogstatsd_non_local_traffic") == true {
		// Listen to all network interfaces
		address = fmt.Sprintf(":%d", config.Datadog.GetInt("dogstatsd_port"))
	} else {
		address = net.JoinHostPort("127.0.0.1", config.Datadog.GetString("dogstatsd_port"))
	}

	url := fmt.Sprintf("udp://%s", address)

	listener := &UDPListener{
		packetPool: packetPool,
		packetBuffer: newPacketBuffer(uint(config.Datadog.GetInt("dogstatsd_packet_buffer_size")),
			config.Datadog.GetDuration("dogstatsd_packet_buffer_flush_timeout"), packetOut),
		url: url,
	}

	evioOptions := evio.Options{
		ReuseInputBuffer: true,
	}
	var events evio.Events
	events.Data = listener.onDatagram
	events.Opened = func(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		return nil, evioOptions, evio.None
	}
	events.Tick = func() (delay time.Duration, action evio.Action) {
		if listener.stopped {
			return time.Second, evio.Shutdown
		}
		return time.Millisecond * 10, evio.None
	}

	listener.events = events
	return listener, nil
}

func (l *UDPListener) onDatagram(c evio.Conn, in []byte) ([]byte, evio.Action) {
	if l.stopped {
		return nil, evio.Shutdown
	}
	packet := l.packetPool.Get()
	copy(packet.buffer, in)
	packet.Contents = packet.buffer[:len(in)]
	// packetBuffer handles the forwarding of the packets to the dogstatsd server intake channel
	l.packetBuffer.append(packet)
	return nil, evio.None
}

func (l *UDPListener) Listen() {
	evio.Serve(l.events, l.url)
}

// Stop closes the UDP connection and stops listening
func (l *UDPListener) Stop() {
	l.stopped = true
	l.packetBuffer.close()
}
