package providers

import (
	"github.com/EngineersBox/Schematic/schema"
	"github.com/EngineersBox/Schematic/state"
)

type Provider struct {
	InstancesMap   map[string]*schema.Instance
	DataSourcesMap map[string]*schema.Data
	ConfigureFunc  ConfigureFunc
}

type ConfigureFunc func(*state.InstanceData) (interface{}, error)
