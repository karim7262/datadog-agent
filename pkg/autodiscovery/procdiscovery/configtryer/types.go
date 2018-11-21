package configtryer

// Config represents a set of informations retrieved
type Config struct {
	Ports       []int    // A list of ports
	UnixSockets []string // A list of paths to unix sockets
}

// ConfigTryer represents a way to retrieve the configuration for the given integration
type ConfigTryer interface {
	Try(string) (*Config, error) // Try to find a config for the provided integration name
}
