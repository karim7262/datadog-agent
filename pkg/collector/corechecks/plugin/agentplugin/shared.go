package agentplugin

import (
	plugin "github.com/hashicorp/go-plugin"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

var PluginMap = map[string]plugin.Plugin{
	"integration": &IntegrationPlugin{},
}

type SenderHelper interface {
	Gauge(metric string, value float64, tags []string) error
	Rate(metric string, value float64, tags []string) error
	Count(metric string, value float64, tags []string) error
	MonotonicCount(metric string, value float64, tags []string) error
	Counter(metric string, value float64, tags []string) error
	Histogram(metric string, value float64, tags []string) error
	Historate(metric string, value float64, tags []string) error
	ServiceCheck(checkName string, status ServiceCheckStatus, tags []string, message string) error
}

type Integration interface {
	Init(name string, initConfig []byte, instances [][]byte) error
	Run(sender SenderHelper, instance []byte) error
}

type IntegrationPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Integration
}

func (p *IntegrationPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterIntegrationServer(s,
		&GRPCServer{
			Impl:   p.Impl,
			broker: broker,
		})
	return nil
}

func (p *IntegrationPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: NewIntegrationClient(c),
		broker: broker,
	}, nil
}

type GRPCClient struct {
	client IntegrationClient
	broker *plugin.GRPCBroker
}

func (m *GRPCClient) Init(name string, initConfig []byte, instances [][]byte) error {
	//TODO: context/timeout
	_, err := m.client.Init(context.Background(), &InitData{
		Name:       name,
		InitConfig: initConfig,
		Instances:  instances,
	})
	return err
}

func (m *GRPCClient) Run(senderHelper SenderHelper, instance []byte) error {
	senderHelperServer := &GRPCSenderHelperServer{Impl: senderHelper}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		RegisterSenderHelperServer(s, senderHelperServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	//TODO: context/timeout
	_, err := m.client.Run(context.Background(),
		&RunData{
			Instance:     instance,
			SenderHandle: brokerID,
		})

	s.Stop()
	return err
}

type GRPCServer struct {
	// This is the real implementation
	Impl   Integration
	broker *plugin.GRPCBroker
}

func (m *GRPCServer) Init(ctx context.Context, req *InitData) (*Empty, error) {
	return &Empty{}, m.Impl.Init(req.Name, req.InitConfig, req.Instances)
}

func (m *GRPCServer) Run(ctx context.Context, req *RunData) (*Empty, error) {
	conn, err := m.broker.Dial(req.SenderHandle)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	sender := &GRPCSenderHelperClient{NewSenderHelperClient(conn)}
	return &Empty{}, m.Impl.Run(sender, req.Instance)
}

type GRPCSenderHelperClient struct {
	client SenderHelperClient
}

func (m *GRPCSenderHelperClient) Gauge(metric string, value float64, tags []string) error {
	_, err := m.client.Gauge(context.Background(), &MetricData{
		Name:  metric,
		Value: value,
		Tags:  tags,
	})
	return err
}

func (m *GRPCSenderHelperClient) Rate(metric string, value float64, tags []string) error {
	_, err := m.client.Rate(context.Background(), &MetricData{
		Name:  metric,
		Value: value,
		Tags:  tags,
	})
	return err
}

func (m *GRPCSenderHelperClient) Count(metric string, value float64, tags []string) error {
	_, err := m.client.Count(context.Background(), &MetricData{
		Name:  metric,
		Value: value,
		Tags:  tags,
	})
	return err
}

func (m *GRPCSenderHelperClient) MonotonicCount(metric string, value float64, tags []string) error {
	_, err := m.client.MonotonicCount(context.Background(), &MetricData{
		Name:  metric,
		Value: value,
		Tags:  tags,
	})
	return err
}

func (m *GRPCSenderHelperClient) Counter(metric string, value float64, tags []string) error {
	_, err := m.client.Counter(context.Background(), &MetricData{
		Name:  metric,
		Value: value,
		Tags:  tags,
	})
	return err
}

func (m *GRPCSenderHelperClient) Histogram(metric string, value float64, tags []string) error {
	_, err := m.client.Histogram(context.Background(), &MetricData{
		Name:  metric,
		Value: value,
		Tags:  tags,
	})
	return err
}

func (m *GRPCSenderHelperClient) Historate(metric string, value float64, tags []string) error {
	_, err := m.client.Historate(context.Background(), &MetricData{
		Name:  metric,
		Value: value,
		Tags:  tags,
	})
	return err
}

func (m *GRPCSenderHelperClient) ServiceCheck(checkName string, status ServiceCheckStatus, tags []string, message string) error {
	_, err := m.client.ServiceCheck(context.Background(), &ServiceCheckData{
		Name:    checkName,
		Status:  ServiceCheckStatus(status),
		Tags:    tags,
		Message: message,
	})
	return err
}

type GRPCSenderHelperServer struct {
	// This is the real implementation
	Impl SenderHelper
}

func (m *GRPCSenderHelperServer) Gauge(ctx context.Context, req *MetricData) (resp *Empty, err error) {
	return &Empty{}, m.Impl.Gauge(req.Name, req.Value, req.Tags)
}

func (m *GRPCSenderHelperServer) Rate(ctx context.Context, req *MetricData) (resp *Empty, err error) {
	return &Empty{}, m.Impl.Rate(req.Name, req.Value, req.Tags)
}

func (m *GRPCSenderHelperServer) Count(ctx context.Context, req *MetricData) (resp *Empty, err error) {
	return &Empty{}, m.Impl.Count(req.Name, req.Value, req.Tags)
}

func (m *GRPCSenderHelperServer) MonotonicCount(ctx context.Context, req *MetricData) (resp *Empty, err error) {
	return &Empty{}, m.Impl.MonotonicCount(req.Name, req.Value, req.Tags)
}

func (m *GRPCSenderHelperServer) Counter(ctx context.Context, req *MetricData) (resp *Empty, err error) {
	return &Empty{}, m.Impl.Counter(req.Name, req.Value, req.Tags)
}

func (m *GRPCSenderHelperServer) Histogram(ctx context.Context, req *MetricData) (resp *Empty, err error) {
	return &Empty{}, m.Impl.Histogram(req.Name, req.Value, req.Tags)
}

func (m *GRPCSenderHelperServer) Historate(ctx context.Context, req *MetricData) (resp *Empty, err error) {
	return &Empty{}, m.Impl.Historate(req.Name, req.Value, req.Tags)
}

func (m *GRPCSenderHelperServer) ServiceCheck(ctx context.Context, req *ServiceCheckData) (resp *Empty, err error) {
	return &Empty{}, m.Impl.ServiceCheck(req.Name, req.Status, req.Tags, req.Message)
}
