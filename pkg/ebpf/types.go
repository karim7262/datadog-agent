package ebpf

// KProbeName stores the name of the kernel probes setup for tracing
type KProbeName string

const (
	// TCPv4DestroySock traces the tcp_v4_destroy_sock system call (called for both ipv4 and ipv6)
	TCPv4DestroySock KProbeName = "kprobe/tcp_v4_destroy_sock"

	// TCPv6Connect traces the v6 connect() system call
	TCPv6Connect KProbeName = "kprobe/tcp_v6_connect"
	// TCPv6ConnectReturn traces the return value for the v6 connect() system call
	TCPv6ConnectReturn KProbeName = "kretprobe/tcp_v6_connect"

	// TCPSendMsg traces the tcp_sendmsg() system call
	TCPSendMsg KProbeName = "kprobe/tcp_sendmsg"

	// TCPSendMsgPre410 traces the tcp_sendmsg() system call on kernels prior to 4.1.0. This is created because
	// we need to load a different kprobe implementation
	TCPSendMsgPre410 KProbeName = "kprobe/tcp_sendmsg/pre_4_1_0"

	// TCPSendMsgReturn traces the return value for the tcp_sendmsg() system call
	// XXX: This is only used for telemetry for now to count the number of errors returned
	// by the tcp_sendmsg func (so we can have a # of tcp sent bytes we miscounted)
	TCPSendMsgReturn KProbeName = "kretprobe/tcp_sendmsg"

	// TCPGetInfo traces the tcp_get_info() system call
	// This probe is used for offset guessing only
	TCPGetInfo KProbeName = "kprobe/tcp_get_info"

	// TCPCleanupRBuf traces the tcp_cleanup_rbuf() system call
	TCPCleanupRBuf KProbeName = "kprobe/tcp_cleanup_rbuf"
	// TCPClose traces the tcp_close() system call
	TCPClose KProbeName = "kprobe/tcp_close"

	// UDPSendMsg traces the udp_sendmsg() system call
	UDPSendMsg KProbeName = "kprobe/udp_sendmsg"
	// UDPSendMsgPre410 traces the udp_sendmsg() system call on kernels prior to 4.1.0
	UDPSendMsgPre410 KProbeName = "kprobe/udp_sendmsg/pre_4_1_0"
	// UDPRecvMsg traces the udp_recvmsg() system call
	UDPRecvMsg KProbeName = "kprobe/udp_recvmsg"
	// UDPRecvMsgPre410 traces the udp_recvmsg() system call on kernels prior to 4.1.0
	UDPRecvMsgPre410 KProbeName = "kprobe/udp_recvmsg/pre_4_1_0"
	// UDPRecvMsgReturn traces the return value for the udp_recvmsg() system call
	UDPRecvMsgReturn KProbeName = "kretprobe/udp_recvmsg"

	// TCPRetransmit traces the return value for the tcp_retransmit_skb() system call
	TCPRetransmit KProbeName = "kprobe/tcp_retransmit_skb"

	// InetCskAcceptReturn traces the return value for the inet_csk_accept syscall
	InetCskAcceptReturn KProbeName = "kretprobe/inet_csk_accept"
)
