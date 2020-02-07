//+build darwin

package inventories

import (
	"bufio"
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"os/exec"
	"strings"
)

//PythonCollector collector for python
type PythonCollector struct {
	pythonPath      string
	environmentName string
}

func (c *PythonCollector) existsInHost() bool {
	cmd := exec.Command("pip", "-V")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (c *PythonCollector) getPackageManagerName() string {
	return fmt.Sprintf("python:%s", c.environmentName)
}

func (c *PythonCollector) collectPackageVersions() *PackageManagerMetadata {
	cmd := exec.Command(c.pythonPath, "-m", "pip", "freeze", "--all")
	out, err := cmd.Output()
	if err != nil {
		return &PackageManagerMetadata{}
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	allPackages := PackageManagerMetadata{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		split := strings.Split(line, "==")
		if len(split) != 2 {
			continue
		}
		packageName, packageVersion := split[0], split[1]

		allPackages[packageName] = packageVersion
	}

	return &allPackages
}

func init() {
	pythonPaths := config.Datadog.GetStringMapString("inventories.python_environments")
	for envName, path := range pythonPaths {
		collector := &PythonCollector{path, envName}
		if collector.existsInHost() {
			RegisterPackageCollector(collector)
		} else {
			log.Warnf("Tried to collect package versions from dpkg. But the agent cannot read data from %s", path)
		}
	}
}
