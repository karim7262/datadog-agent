package inventories

//PackageCollector is a mapping from pkg mngr name to metadata.
type PackageCollector interface {
	existsInHost() bool
	getPackageManagerName() string
	collectPackageVersions() *PackageManagerMetadata
}

var (
	//PackageCollectors is a mapping from pkg mngr name to metadata.
	PackageCollectors = make(map[string]PackageCollector)
)

//RegisterPackageCollector abc
func RegisterPackageCollector(collector PackageCollector) {
	name := collector.getPackageManagerName()
	PackageCollectors[name] = collector
}

//CollectPackagesVersions abc
func CollectPackagesVersions() *PackagesMetadata {
	packagesMetadata := make(PackagesMetadata)
	for collectorName, collector := range PackageCollectors {
		packagesMetadata[collectorName] = collector.collectPackageVersions()
	}

	return &packagesMetadata
}
