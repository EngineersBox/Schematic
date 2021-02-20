package schema

import (
	"fmt"
	"sync"
)

type InstanceState struct {
	// ID is the name gives to the instance
	ID         string                 `json:"id"`
	Type       string                 `json:"kind"`
	Attributes map[string]interface{} `json:"attributes"`
	Meta       map[string]interface{} `json:"meta"`
	Provider   string                 `json:"provider"`
	// Tainted is used to mark a resource for recreation.
	Tainted bool `json:"tainted"`

	mu sync.Mutex
}

func (s *InstanceState) Lock()   { s.mu.Lock() }
func (s *InstanceState) Unlock() { s.mu.Unlock() }

func (s *InstanceState) init() {
	s.Lock()
	defer s.Unlock()

	if s.Attributes == nil {
		s.Attributes = make(map[string]interface{})
	}
	if s.Meta == nil {
		s.Meta = make(map[string]interface{})
	}
}

func GetInstanceAttributeNesting(attributeNesting []string, attributes map[string]interface{}) (interface{}, error) {
	if len(attributeNesting) == 0 {
		return nil, fmt.Errorf("must have at least 1 field specified")
	}
	if len(attributeNesting) == 1 {
		bottomNesting, ok := attributes[attributeNesting[0]]
		if !ok {
			return nil, fmt.Errorf("no such field: %s", attributeNesting[0])
		}
		return bottomNesting, nil
	}
	nestedField, ok := attributes[attributeNesting[0]].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no such field: %s", attributeNesting[0])
	}
	return GetInstanceAttributeNesting(attributeNesting[1:], nestedField)
}
