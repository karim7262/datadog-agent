package providers

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/DataDog/datadog-agent/pkg/collector/check"
	"github.com/DataDog/datadog-agent/pkg/config"
	log "github.com/cihub/seelog"
	"github.com/ghodss/yaml"
)

const (
	instancePath   string = "instances"
	checkNamePath  string = "check_names"
	initConfigPath string = "init_configs"
)

// parseJSONValue returns a slice of ConfigData parsed from the JSON
// contained in the `value` parameter
func parseJSONValue(value string) ([]check.ConfigData, error) {
	if value == "" {
		return nil, fmt.Errorf("Value is empty")
	}

	yamlValue, err := yaml.JSONToYAML([]byte(value))
	if err != nil {
		return nil, err
	}

	var rawRes interface{}
	var result []check.ConfigData

	err = yaml.Unmarshal(yamlValue, &rawRes)
	if err != nil {
		return nil, err
	}

	for _, r := range rawRes.([]interface{}) {
		switch r.(type) {
		case []byte:
			result = append(result, r.([]byte))
		case map[string]interface{}:
			init, _ := yaml.Marshal(r)
			result = append(result, init)
		}

	}
	return result, nil
}

func parseCheckNames(names string) (res []string, err error) {
	if names == "" {
		return nil, fmt.Errorf("check_names is empty")
	}

	if err = json.Unmarshal([]byte(names), &res); err != nil {
		return nil, err
	}

	return res, nil
}

func buildStoreKey(key ...string) string {
	parts := []string{config.Datadog.GetString("autoconf_template_dir")}
	parts = append(parts, key...)
	return path.Join(parts...)
}

func buildTemplates(key string, checkNames []string, initConfigs, instances []check.ConfigData) []check.Config {
	templates := make([]check.Config, 0)

	// sanity check
	if len(checkNames) != len(initConfigs) || len(checkNames) != len(instances) {
		log.Error("Template entries don't all have the same length in etcd, not using them.")
		return templates
	}

	for idx := range checkNames {
		instance := check.ConfigData(instances[idx])

		templates = append(templates, check.Config{
			ID:         check.ID(key),
			Name:       checkNames[idx],
			InitConfig: check.ConfigData(initConfigs[idx]),
			Instances:  []check.ConfigData{instance},
		})
	}
	return templates
}

func buildConfigs(id string, templates []check.Config) []check.Config {
	configs := make([]check.Config, 0)
	for _, template := range templates {
		c := check.Config{
			ID:         check.ID(id),
			Name:       template.Name,
			InitConfig: template.InitConfig,
			Instances:  template.Instances,
		}
		configs = append(configs, c)
	}

	return configs
}
