package schema

import (
	"os"

	"github.com/EngineersBox/Schematic/schematic"
)

type Schema struct {
	Type ValueType

	Optional bool
	Required bool

	DiffFunc SchemaDiffFunc

	Description string

	Computed  bool
	ForceNew  bool
	StateFunc SchemaStateFunc

	// The following fields are only set for a TypeList, TypeSet, or TypeMap.
	//
	// Elem represents the element type. For a TypeMap, it must be a *Schema
	// with a Type that is one of the primitives: TypeString, TypeBool,
	// TypeInt, or TypeFloat. Otherwise it may be either a *Schema or a
	// *Resource. If it is *Schema, the element type is just a simple value.
	// If it is *Resource, the element type is a complex structure,
	// potentially managed via its own CRUD actions on the API.
	Elem interface{}

	// The following fields are only set for a TypeList or TypeSet.
	//
	// MaxItems defines a maximum amount of items that can exist within a
	// TypeSet or TypeList. Specific use cases would be if a TypeSet is being
	// used to wrap a complex structure, however more than one instance would
	// cause instability.
	//
	// MinItems defines a minimum amount of items that can exist within a
	// TypeSet or TypeList. Specific use cases would be if a TypeSet is being
	// used to wrap a complex structure, however less than one instance would
	// cause instability.
	//
	// If the field Optional is set to true then MinItems is ignored and thus
	// effectively zero.
	MaxItems int
	MinItems int

	// The following fields are only valid for a TypeSet type.
	//
	// Set defines a function to determine the unique ID of an item so that
	// a proper set can be built.
	Set SchemaSetFunc

	// ConflictsWith is a set of schema keys that conflict with this schema.
	// This will only check that they're set in the _config_. This will not
	// raise an error for a malfunctioning resource that sets a conflicting
	// key.
	//
	// ExactlyOneOf is a set of schema keys that, when set, only one of the
	// keys in that list can be specified. It will error if none are
	// specified as well.
	//
	// AtLeastOneOf is a set of schema keys that, when set, at least one of
	// the keys in that list must be specified.
	//
	// RequiredWith is a set of schema keys that must be set simultaneously.
	ConflictsWith []string
	ExactlyOneOf  []string
	AtLeastOneOf  []string
	RequiredWith  []string

	Deprecated bool

	ValidateFunc SchemaValidateFunc

	// Sensitive ensures that the attribute's value does not get displayed in
	// logs or regular output. It should be used for passwords or other
	// secret fields. Future versions of Terraform may encrypt these
	// values.
	Sensitive bool
}

type SchemaDefaultFunc func() (interface{}, error)

type SchemaSetFunc func(interface{}) int

type SchemaStateFunc func(interface{}) string

type SchemaValidateFunc func(interface{}, string) ([]string, []error)

type SchemaDiffFunc func(k, old, new string, d *InstanceData) bool

type schemaMap map[string]*Schema

func (m schemaMap) panicOnError() bool {
	return os.Getenv("TF_ACC") != ""
}

// Data returns a ResourceData for the given schema, state, and diff.
//
// The diff is optional.
func (m schemaMap) Data(
	s *schematic.InstanceState,
	d *schematic.InstanceDiff) (*InstanceData, error) {
	return &InstanceData{
		schema:       m,
		state:        s,
		diff:         d,
		panicOnError: m.panicOnError(),
	}, nil
}
