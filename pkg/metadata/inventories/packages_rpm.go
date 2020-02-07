package inventories

import (
	"bufio"
	"errors"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
)

var (
	ignoredPackages = []string{
		"gpg-pubkey", // Ignore gpg-pubkey packages which are fake packages used to store GPG keys - they are not versionned properly.
	}
	allowedSymbols = []rune{'.', '-', '+', '~', ':', '_'}
)

//RpmCollector collector for rpm
type RpmCollector struct {
	rpmBin string
}

func (c *RpmCollector) existsInHost() bool {
	cmd := exec.Command(c.rpmBin, "--version")
	if err := cmd.Run(); err != nil {
		//fmt.Println(string(out))
		return false
	}
	return true
}

func (c *RpmCollector) getPackageManagerName() string {
	return "rpm"
}

func (c *RpmCollector) collectPackageVersions() *PackageManagerMetadata {
	allPackages := PackageManagerMetadata{}
	out, err := exec.Command(c.rpmBin, "-qa", "--qf", "%{NAME} %{EPOCH}:%{VERSION}-%{RELEASE} %{SOURCERPM}\n").CombinedOutput()
	if err != nil {
		return &allPackages
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		packageName, packageVersion := parseRpmOutput(scanner.Text())
		allPackages[packageName] = packageVersion
	}
	return &allPackages
}

func parseRpmOutput(raw string) (packageName string, version string) {
	line := strings.Split(raw, " ")
	if len(line) != 3 {
		// We may see warnings on some RPM versions:
		// "warning: Generating 12 missing index(es), please wait..."
		return
	}

	if isIgnored(line[0]) {
		return
	}

	packageName, version = line[0], strings.Replace(line[1], "(none):", "", -1)
	if err := validateVersion(version); err != nil {
		log.Error(err)
		return
	}
	return
}

func isIgnored(packageName string) bool {
	for _, pkg := range ignoredPackages {
		if pkg == packageName {
			return true
		}
	}

	return false
}

func validateVersion(versionStr string) error {
	// Trim leading and trailing space
	str := strings.TrimSpace(versionStr)
	var version string
	var release string

	if len(str) == 0 {
		return errors.New("Version string is empty")
	}

	// Find epoch
	sepepoch := strings.Index(str, ":")
	if sepepoch > -1 {
		intepoch, err := strconv.Atoi(str[:sepepoch])
		if err != nil {
			return errors.New("epoch in version is not a number")
		}
		if intepoch < 0 {
			return errors.New("epoch in version is negative")
		}
	}

	// Find version / release
	seprevision := strings.Index(str, "-")
	if seprevision > -1 {
		version = str[sepepoch+1 : seprevision]
		release = str[seprevision+1:]
	} else {
		version = str[sepepoch+1:]
		release = ""
	}
	// Verify format
	if len(version) == 0 {
		return errors.New("No version")
	}

	for i := 0; i < len(version); i = i + 1 {
		r := rune(version[i])
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !validSymbol(r) {
			return errors.New("invalid character in version")
		}
	}

	for i := 0; i < len(release); i = i + 1 {
		r := rune(release[i])
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !validSymbol(r) {
			return errors.New("invalid character in revision")
		}
	}

	return nil
}

func validSymbol(r rune) bool {
	return containsRune(allowedSymbols, r)
}

func containsRune(s []rune, e rune) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func init() {
	rpmBin := config.Datadog.GetString("inventories.rpm_bim")
	collector := &RpmCollector{rpmBin}
	if collector.existsInHost() {
		RegisterPackageCollector(collector)
	} else {
		log.Warnf("Tried to collect package versions from rpm. But the agent cannot use the binary %s", rpmBin)
	}
}
