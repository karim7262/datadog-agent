package network

import (
	"path"
	"time"

	"github.com/DataDog/datadog-agent/pkg/collector/check"
	core "github.com/DataDog/datadog-agent/pkg/collector/corechecks"
	"gopkg.in/yaml.v2"
)

const (
	netDevPath = "/proc/1/net/dev"
)

func init() {
	core.RegisterCheck("netdev", newNetDevCheck)
}

func newNetDevCheck() check.Check {
	cfg := &netdevConfig{}
	return &NetDevCheck{
		cfg: cfg,
	}
}

type NetDevCheck struct {
	id           check.ID
	lastWarnings []error
	cfg          *netdevConfig
}

type netdevInstanceConfig struct {
	ProcPrefix string `yaml:"proc_prefix"`
	NetDevPath string
}

type netdevInitConfig struct{}

type netdevConfig struct {
	instance netdevInstanceConfig
	initConf netdevInitConfig
}

// Run the check
func (c *NetDevCheck) Run() error {
	return nil
}

// Stop the check if it's running
func (c *NetDevCheck) Stop() {

}

// String provide a printable version of the check name
func (c *NetDevCheck) String() string {
	return "net/dev"
}

// Configure the check from the outside
func (c *NetDevCheck) Configure(config, initConfig check.ConfigData) error {
	err := c.parse(config, initConfig)
	if err != nil {
		return err
	}
	c.id = check.Identify(c, config, initConfig)
	return nil
}

// Interval return the interval time for the check
func (c *NetDevCheck) Interval() time.Duration {
	return check.DefaultCheckInterval
}

// ID provide a unique identifier for every check instance
func (c *NetDevCheck) ID() check.ID {
	return c.id
}

// GetWarnings return the last warning registered by the check
func (c *NetDevCheck) GetWarnings() []error {
	return nil
}

// GetMetricStats get metric stats from the sender
func (c *NetDevCheck) GetMetricStats() (map[string]int64, error) {
	return nil, nil
}

// parse
func (c *NetDevCheck) parse(instanceData, initData []byte) error {
	err := yaml.Unmarshal(instanceData, &c.cfg.instance)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(initData, &c.cfg.initConf)
	if err != nil {
		return err
	}

	c.cfg.instance.NetDevPath = path.Join(c.cfg.instance.ProcPrefix, netDevPath)
	return nil
}

func (c *NetDevCheck) getNetDevPath() string {
	return c.cfg.instance.NetDevPath
}
