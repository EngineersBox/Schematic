package state

import (
	"github.com/EngineersBox/Schematic/schema"
	"github.com/zclconf/go-cty/cty"
)

type Variable struct {
	Name     string
	Value    cty.Value `json:"Value"`
	BaseType schema.ValueType
}
