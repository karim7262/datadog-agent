package plugin

import (
	"time"

	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/plugin/agentplugin"
	"github.com/DataDog/datadog-agent/pkg/metrics"
)

type PluginCheck struct {
	runCommand string
	name       string
	plugin     agentplugin.Integration
	instance   check.ConfigData
	id         check.ID
}

func NewPluginCheck(name string, integration agentplugin.Integration) *PluginCheck {
	return &PluginCheck{
		name:   name,
		plugin: integration,
	}
}

func (c *PluginCheck) Run() error {
	sender, err := aggregator.GetSender(c.ID())
	if err != nil {
		return err
	}

	proxy := &SenderProxy{
		sender: sender,
	}

	err = c.plugin.Run(proxy, c.instance)
	sender.Commit()

	return err
}

func (c *PluginCheck) Stop() {

}

func (c *PluginCheck) String() string {
	return c.name
}

func (c *PluginCheck) Configure(config check.ConfigData, initConfig check.ConfigData) error {
	c.instance = config
	c.id = check.Identify(c, config, initConfig)
	return nil
}

func (c *PluginCheck) Interval() time.Duration {
	return 15 * time.Second
}

//FIXME: Identify.BuildID
func (c *PluginCheck) ID() check.ID {
	return c.id
}

func (c *PluginCheck) GetWarnings() []error {
	return []error{}
}

func (c *PluginCheck) GetMetricStats() (map[string]int64, error) {
	return map[string]int64{}, nil
}

type SenderProxy struct {
	sender aggregator.Sender
}

func (s *SenderProxy) Gauge(metric string, value float64, tags []string) error {
	s.sender.Gauge(metric, value, "", tags)
	return nil
}

func (s *SenderProxy) Rate(metric string, value float64, tags []string) error {
	s.sender.Rate(metric, value, "", tags)
	return nil
}

func (s *SenderProxy) Count(metric string, value float64, tags []string) error {
	s.sender.Count(metric, value, "", tags)
	return nil
}

func (s *SenderProxy) MonotonicCount(metric string, value float64, tags []string) error {
	s.sender.MonotonicCount(metric, value, "", tags)
	return nil
}

func (s *SenderProxy) Counter(metric string, value float64, tags []string) error {
	s.sender.Counter(metric, value, "", tags)
	return nil
}

func (s *SenderProxy) Histogram(metric string, value float64, tags []string) error {
	s.sender.Histogram(metric, value, "", tags)
	return nil
}

func (s *SenderProxy) Historate(metric string, value float64, tags []string) error {
	s.sender.Historate(metric, value, "", tags)
	return nil
}

func (s *SenderProxy) ServiceCheck(checkName string, status agentplugin.ServiceCheckStatus, tags []string, message string) error {
	s.sender.ServiceCheck(checkName, metrics.ServiceCheckStatus(status), "", tags, message)
	return nil
}
