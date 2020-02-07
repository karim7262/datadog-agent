//+build Linux

package inventories

import (
	"bufio"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	versionAllowedSymbols     = []rune{'.', '-', '+', '~', ':', '_'}
	revisionAllowedSymbols    = []rune{'.', '+', '~', '_'}
	dpkgSrcCaptureRegexp      = regexp.MustCompile(`Source: (?P<name>[^\s]*)( \((?P<version>.*)\))?`)
	dpkgSrcCaptureRegexpNames = dpkgSrcCaptureRegexp.SubexpNames()
)

type DpkgCollector struct {
	dpkgStatusFile string
}

func (c *DpkgCollector) existsInHost() bool {
	_, err := os.Open(dpkgStatusFile)
	return err == nil
}

func (c *DpkgCollector) getPackageManagerName() string {
	return "dpkg"
}

func (c *DpkgCollector) collectPackageVersions() *PackageManagerMetadata {
	fileHandle, err := os.Open(dpkgStatusFile)
	if err != nil {
		return &PackageManagerMetadata{}
	}
	defer fileHandle.Close()
	scanner := bufio.NewScanner(fileHandle)
	allPackages := make(PackageManagerMetadata)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		packageName, version := parseDpkgStatus(scanner)
		allPackages[packageName] = version
	}

	return &allPackages
}

func parseDpkgStatus(scanner *bufio.Scanner) (packageName string, version string) {
	var (
		name          string
		sourceName    string
		sourceVersion string
	)
	for {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			// End of package description
			break
		}

		if strings.HasPrefix(line, "Package: ") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "Package: "))
		} else if strings.HasPrefix(line, "Source: ") {
			// Source line (Optional)
			// Gives the name of the source package
			// May also specifies a version
			srcCapture := dpkgSrcCaptureRegexp.FindAllStringSubmatch(line, -1)[0]
			md := map[string]string{}
			for i, n := range srcCapture {
				md[dpkgSrcCaptureRegexpNames[i]] = strings.TrimSpace(n)
			}

			sourceName = md["name"]
			if md["version"] != "" {
				sourceVersion = md["version"]
			}
		} else if strings.HasPrefix(line, "Version: ") {
			version = strings.TrimPrefix(line, "Version: ")
		}

		if !scanner.Scan() {
			break
		}
	}

	if name != "" && version != "" {
		if isValidVersion(version) {
			return name, version
		} else {
			log.Debugf("Unable to validate version %s for package %s.", name, version)
		}
	}

	if sourceName != "" && sourceVersion != "" {
		sourceVersion = version
	} else if sourceName != "" && sourceVersion == "" {
		sourceName = name
	}

	if sourceName != "" && sourceVersion != "" {
		if isValidVersion(sourceVersion) {
			return sourceName, sourceVersion
		} else {
			log.Debugf("Unable to validate version %s for package %s.", sourceName, sourceVersion)
		}
	}
	return
}

func isValidVersion(version string) bool {
	var (
		packageVersion string
		revision       string
	)
	version = strings.TrimSpace(version)

	if len(version) == 0 {
		return false
	}

	// Find epoch
	sepepoch := strings.Index(version, ":")
	if sepepoch > -1 {
		intepoch, err := strconv.Atoi(version[:sepepoch])
		if err != nil {
			return false
		}
		if intepoch < 0 {
			return false
		}
	}

	// Find version / revision
	seprevision := strings.LastIndex(version, "-")
	if seprevision > -1 {
		packageVersion = version[sepepoch+1 : seprevision]
		revision = version[seprevision+1:]
	} else {
		packageVersion = version[sepepoch+1:]
		revision = ""
	}
	// Verify format
	if len(packageVersion) == 0 {
		return false
	}

	for i := 0; i < len(packageVersion); i = i + 1 {
		r := rune(packageVersion[i])
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !containsRune(versionAllowedSymbols, r) {
			return false
		}
	}

	for i := 0; i < len(revision); i = i + 1 {
		r := rune(revision[i])
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) && !containsRune(revisionAllowedSymbols, r) {
			return false
		}
	}

	return true
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
	dpkgStatusFile := config.Datadog.GetString("inventories.dpkg_status_file")
	collector := &DpkgCollector{dpkgStatusFile}
	if collector.existsInHost() {
		RegisterPackageCollector(collector)
	} else {
		log.Warnf("Tried to collect package versions from dpkg. But the agent cannot read data from %s", dpkgStatusFile)
	}
}
