package schema

/*
 * This file is basically copied from the HashiCorp Terraform source code:
 * https://github.com/hashicorp/terraform-plugin-sdk/blob/master/helper/schema/resource_data.go
 */

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/EngineersBox/Schematic/schematic"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type InstanceData struct {
	// Settable (internally)
	schema       map[string]*Schema
	config       *schematic.InstanceConfig
	state        *schematic.InstanceState
	diff         *schematic.InstanceDiff
	meta         map[string]interface{}
	timeouts     *InstanceTimeout
	providerMeta cty.Value

	// multiReader *MultiLevelFieldReader
	// setWriter   *MapFieldWriter
	newState *schematic.InstanceState
	partial  bool
	once     sync.Once
	isNew    bool

	panicOnError bool
}

// getResult is the internal structure that is generated when a Get
// is called that contains some extra data that might be used.
type getResult struct {
	Value          interface{}
	ValueProcessed interface{}
	Computed       bool
	Exists         bool
	Schema         *Schema
}

func (d *InstanceData) Get(key string) interface{} {
	return nil // TODO: implement this
}

func (d *InstanceData) GetChange(key string) (interface{}, interface{}) {
	return nil, nil // TODO: implement this
}

func (d *InstanceData) HasChanges(keys ...string) bool {
	for _, key := range keys {
		if d.HasChange(key) {
			return true
		}
	}
	return false
}

func (d *InstanceData) HasChangesExcept(keys ...string) bool {
	for attr := range d.diff.Attributes {
		rootAttr := strings.Split(attr, ".")[0]
		var skipAttr bool

		for _, key := range keys {
			if rootAttr == key {
				skipAttr = true
				break
			}
		}

		if !skipAttr && d.HasChange(rootAttr) {
			return true
		}
	}

	return false
}

func (d *InstanceData) HasChange(key string) bool {
	o, n := d.GetChange(key)

	if eq, ok := o.(Equal); ok {
		return !eq.Equal(n)
	}

	return !reflect.DeepEqual(o, n)
}

func (d *InstanceData) HasChangeExcept(key string) bool {
	for attr := range d.diff.Attributes {
		rootAttr := strings.Split(attr, ".")[0]

		if rootAttr == key {
			continue
		}

		if d.HasChange(rootAttr) {
			return true
		}
	}

	return false
}

func (d *InstanceData) Set(key string, value interface{}) error {
	return nil // TODO: implement this
}

func (d *InstanceData) MarkNewResource() {
	d.isNew = true
}

func (d *InstanceData) IsNewResource() bool {
	return d.isNew
}

func (d *InstanceData) Id() string {
	var result string

	if d.state != nil {
		result = d.state.ID
		if result == "" {
			result = d.state.Attributes["id"]
		}
	}

	if d.newState != nil {
		result = d.newState.ID
		if result == "" {
			result = d.newState.Attributes["id"]
		}
	}

	return result
}

func (d *InstanceData) SetId(v string) {
	// TODO: implement this
}

func (d *InstanceData) SetType(t string) {
	// TODO: implement this
}

// State returns the new InstanceState after the diff and any Set
// calls.
func (d *InstanceData) State() *schematic.InstanceState {
	return nil // TODO: implement this
}

func (d *InstanceData) Timeout(key string) time.Duration {
	key = strings.ToLower(key)

	// System default of 20 minutes
	defaultTimeout := 20 * time.Minute

	if d.timeouts == nil {
		return defaultTimeout
	}

	var timeout *time.Duration
	switch key {
	case TimeoutCreate:
		timeout = d.timeouts.Create
	case TimeoutRead:
		timeout = d.timeouts.Read
	case TimeoutUpdate:
		timeout = d.timeouts.Update
	case TimeoutDelete:
		timeout = d.timeouts.Delete
	}

	if timeout != nil {
		return *timeout
	}

	if d.timeouts.Default != nil {
		return *d.timeouts.Default
	}

	return defaultTimeout
}

func (d *InstanceData) init() {
	// TODO: implement this
}

func (d *InstanceData) GetProviderMeta(dst interface{}) error {
	if d.providerMeta.IsNull() {
		return nil
	}
	return gocty.FromCtyValue(d.providerMeta, &dst)
}
