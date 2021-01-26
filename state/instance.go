package state

import "fmt"

type Instance struct {
	Provider string                 `json:"provider"`
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	Fields   map[string]interface{} `json:"fields"`
}

func GetInstanceField(fieldNesting []string, fields map[string]interface{}) (interface{}, error) {
	if len(fieldNesting) == 0 {
		return nil, fmt.Errorf("must have at least 1 field specified")
	}
	if len(fieldNesting) == 1 {
		return fields[fieldNesting[0]], nil
	}
	nestedField, ok := fields[fieldNesting[0]].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("not such field: %s", fieldNesting[0])
	}
	return GetInstanceField(fieldNesting[1:], nestedField)
}
