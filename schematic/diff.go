package schematic

import (
	"regexp"
	"sync"
)

type diffChangeType byte

const (
	diffInvalid diffChangeType = iota
	diffNone
	diffCreate
	diffUpdate
	diffDestroy
	diffDestroyCreate
)

// multiVal matches the index key to a flatmapped set, list or map
var multiVal = regexp.MustCompile(`\.(#|%)$`)

type InstanceDiff struct {
	mu sync.Mutex

	Attributes     map[string]*InstanceAttrDiff // TODO: Add this
	Destroy        bool
	DestroyDeposed bool
	DestroyTainted bool

	Meta map[string]interface{}
}

func (d *InstanceDiff) Lock()   { d.mu.Lock() }
func (d *InstanceDiff) Unlock() { d.mu.Unlock() }

type InstanceAttrDiff struct {
	Old         string      // Old Value
	New         string      // New Value
	NewComputed bool        // True if new value is computed (unknown currently)
	NewRemoved  bool        // True if this attribute is being removed
	NewExtra    interface{} // Extra information for the provider
	RequiresNew bool        // True if change requires new resource
	Sensitive   bool        // True if the data should not be displayed in UI output
}
