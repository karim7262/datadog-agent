package agent

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	tracerBcastPort = 1911
	agentBcastPort  = 1912
)

const aliveMessage = `{"status":"alive","source":"trace-agent","port":%q}`

func runDiscovery(ctx context.Context, src, port string) {
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(src, strconv.Itoa(tracerBcastPort)))
	if err != nil {
		log.Errorf("Could not set up broadcast. Error resolving address: %v", err)
		return
	}
	udpc, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Errorf("Could not set up broadcast. Error listening on UDP: %v", err)
		return
	}
	defer udpc.Close()

	msg := fmt.Sprintf(aliveMessage, port)
	go listenBroadcast(ctx, udpc, msg)

	for i := 0; i < 2; i++ {
		// announce two times at startup, then take it easy
		select {
		case <-ctx.Done():
			return
		default:
			// continue
		}
		doBroadcast(udpc, msg)
		time.Sleep(time.Second)
	}
	for {
		// keep announcing every 10 seconds
		select {
		case <-ctx.Done():
			return
		default:
			// continue
		}
		doBroadcast(udpc, msg)
		time.Sleep(2 * time.Second)
	}
}

func doBroadcast(conn *net.UDPConn, msg string) {
	_, err := conn.WriteTo([]byte(msg), &net.UDPAddr{IP: net.IPv4bcast, Port: agentBcastPort})
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return
		}
		log.Debugf("Error broadcasting presence: %v", err)
	}
}

func listenBroadcast(ctx context.Context, conn *net.UDPConn, msg string) {
	in := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// continue
		}
		conn.SetDeadline(time.Now().Add(2 * time.Second))
		_, addr, err := conn.ReadFrom(in)
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			continue
		}
		if err != nil {
			log.Debugf("Error reading request: %v", err)
			continue
		}
		_, err = conn.WriteTo([]byte(msg), addr)
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			continue
		}
		if err != nil {
			log.Debugf("Error reading request: %v", err)
			continue
		}
	}
}
