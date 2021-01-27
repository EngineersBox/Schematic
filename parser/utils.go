package parser

import (
	"fmt"
	"github.com/EngineersBox/Schematic/schema"
)

func validateSchemaFields(fields []string, instanceSchema map[string]*schema.Schema) bool {
	hasField := instanceSchema[fields[0]] != nil
	if len(fields) == 1 {
		return hasField
	}
	if instanceSchema[fields[0]].Type != schema.TypeMap {
		return false
	}
	nestedSchema, ok := instanceSchema[fields[0]].Elem.(map[string]*schema.Schema)
	if !ok {
		return false
	}
	return validateSchemaFields(fields[1:], nestedSchema)
}

func recurseAssign(nesting []string, value interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	if len(nesting) == 1 {
		obj[nesting[0]] = value
		return obj, nil
	}
	next := nesting[0]
	if obj[next] == nil {
		obj[next] = make(map[string]interface{})
	}
	nextMap, ok := obj[next].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not update nesting")
	}
	recusedValue, err := recurseAssign(nesting[1:], value, nextMap)
	obj[next] = recusedValue
	if err != nil {
		return nil, err
	}
	return obj, nil
}
