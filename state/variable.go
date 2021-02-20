package state

import (
	"github.com/EngineersBox/Schematic/collection"
	"github.com/zclconf/go-cty/cty"
)

type Variable struct {
	Name     string
	Value    cty.Value `json:"Value"`
	BaseType schematic.ValueType
}
