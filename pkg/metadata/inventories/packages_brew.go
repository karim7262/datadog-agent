//+build darwin

package inventories

import (
	"bufio"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"os/exec"
	"strings"
)

type BrewCollector struct {
	brewPath string
}

func (c *BrewCollector) existsInHost() bool {
	cmd := exec.Command(c.brewPath, "-v")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (c *BrewCollector) getPackageManagerName() string {
	return "brew"
}

func (c *BrewCollector) collectPackageVersions() *PackageManagerMetadata {
	cmd := exec.Command(c.brewPath, "list", "--versions")
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
		split := strings.Split(line, " ")
		if len(split) < 2 {
			continue
		}
		packageName, packageVersion := split[0], split[1]

		allPackages[packageName] = packageVersion
	}

	return &allPackages
}

func init() {
	brewPath := config.Datadog.GetString("inventories.brew_path")
	collector := &BrewCollector{brewPath}
	if collector.existsInHost() {
		RegisterPackageCollector(collector)
	} else {
		log.Warnf("Tried to collect package versions from brew. But brew does not seem installed.")
	}
}
