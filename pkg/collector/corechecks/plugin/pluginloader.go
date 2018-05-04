package plugin

import (
	"fmt"
	"os/exec"

	"github.com/DataDog/datadog-agent/pkg/collector/check"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/plugin/agentplugin"
	"github.com/DataDog/datadog-agent/pkg/collector/loaders"
	plugin "github.com/hashicorp/go-plugin"
)

type PluginCheckLoader struct {
	checks []string
}

func NewPluginCheckLoader() *PluginCheckLoader {
	return &PluginCheckLoader{
		checks: []string{},
	}
}

func (pl *PluginCheckLoader) Load(config check.Config) ([]check.Check, error) {
	if config.Plugin != "" {
		checks := []check.Check{}

		cmd := exec.Command("sh", "-c", config.Plugin)
		cmd.Dir = "/opt/datadog-agent/etc/plugins"
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig:  agentplugin.Handshake,
			Plugins:          agentplugin.PluginMap,
			Cmd:              cmd,
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			Managed:          true,
		})

		rpcClient, err := client.Client()
		if err != nil {
			return nil, err
		}

		raw, err := rpcClient.Dispense("integration")
		if err != nil {
			return nil, err
		}

		integration := raw.(agentplugin.Integration)

		integration.Init(config.Name, config.InitConfig, [][]byte{})

		for _, i := range config.Instances {
			check := NewPluginCheck(config.Name, integration)
			check.Configure(i, config.InitConfig)
			checks = append(checks, check)
		}

		fmt.Println("Loaded", config.Name, len(checks))
		return checks, nil
	}

	return nil, fmt.Errorf("not a gRPC plugin check")
}

func (pl *PluginCheckLoader) String() string {
	return "gRPC Plugin Check Loader"
}

func init() {
	factory := func() (check.Loader, error) {
		return NewPluginCheckLoader(), nil
	}

	loaders.RegisterLoader(1337, factory)
}
