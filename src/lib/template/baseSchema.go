package schematic

import (
	"fmt"
	"log"
)

type Schema struct {
	Resources []Resource
}

type Resource struct {
	Name       string
	Type       string
	Attributes []Attribute
}

type ResourceError struct {
	arg  int
	prob string
}

func (e *ResourceError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

type Attribute struct {
	Key   string
	Value interface{}
}

type AttributeError struct {
	arg  int
	prob string
}

func (e *AttributeError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.prob)
}

// BuildSchema ... Construct a schema object from a generic interface, erroring when the structure is invalid
func (s *Schema) BuildSchema(resArr []interface{}) {
	s.Resources = make([]Resource, len(resArr))
	for _, res := range resArr {
		if err := ValidateResource(res); err != nil {
			log.Fatalln(err)
		}
		s.Resources = append(s.Resources, res.(Resource))
	}
}

// ValidateResource ... Assert the structure of a Resource to a generic interface, erroring when invalid
func ValidateResource(res interface{}) error {
	resolvedRes, valid := res.(Resource)
	if !valid {
		return &ResourceError{1, "Invalid resource"}
	}
	for _, attr := range resolvedRes.Attributes {
		if err := ValidateAttribute(attr); err != nil {
			return err
		}
	}
	return nil
}

// ValidateAttribute ... Assert the structure of an Attribute to a generic interface, erroring when invalid
func ValidateAttribute(attr interface{}) error {
	_, valid := attr.(Attribute)
	if !valid {
		return &AttributeError{1, "Invalid attribute"}
	}
	return nil
}
