package inventories

type PackageCollector interface {
	existsInHost() bool
	getPackageManagerName() string
	collectPackageVersions() *PackageManagerMetadata
}

var (
	PackageCollectors = make(map[string]PackageCollector)
)

func RegisterPackageCollector(collector PackageCollector) {
	name := collector.getPackageManagerName()
	PackageCollectors[name] = collector
}

func CollectPackagesVersions() *PackagesMetadata {
	packagesMetadata := make(PackagesMetadata)
	for collectorName, collector := range PackageCollectors {
		packagesMetadata[collectorName] = collector.collectPackageVersions()
	}

	return &packagesMetadata
}
