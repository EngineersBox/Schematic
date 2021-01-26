package schema

type Provider struct {
	Schema         map[string]*Schema
	InstancesMap   map[string]*Instance
	DataSourcesMap map[string]*Data
	ConfigureFunc  ConfigureFunc
}

type ConfigureFunc func(*InstanceData) (interface{}, error)
