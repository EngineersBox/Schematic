package schematic

import (
	"sync"

	"github.com/zclconf/go-cty/cty"
)

type InstanceState struct {
	ID           string                 `json:"id"`
	Attributes   map[string]string      `json:"attributes"`
	Meta         map[string]interface{} `json:"meta"`
	ProviderMeta cty.Value
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
		s.Attributes = make(map[string]string)
	}
	if s.Meta == nil {
		s.Meta = make(map[string]interface{})
	}
}
