package schema

import (
	"context"

	"github.com/EngineersBox/Schematic/schematic"
)

var ReservedDataSourceFields = []string{
	"connection",
	"count",
	"depends_on",
	"lifecycle",
	"provider",
	"provisioner",
}

var ReservedResourceFields = []string{
	"connection",
	"count",
	"depends_on",
	"lifecycle",
	"provider",
	"provisioner",
}

type Instance struct {
	Schema map[string]*Schema

	Create CreateFunc
	Read   ReadFunc
	Update UpdateFunc
	Delete DeleteFunc

	Exists ExistsFunc

	CustomizeDiff CustomizeDiffFunc

	DeprecationMessage string

	Timeouts *InstanceTimeout

	Description string
}

type CreateFunc func(*InstanceData, interface{}) error
type ReadFunc func(*InstanceData, interface{}) error
type UpdateFunc func(*InstanceData, interface{}) error
type DeleteFunc func(*InstanceData, interface{}) error
type ExistsFunc func(*InstanceData, interface{}) (bool, error)

type CustomizeDiffFunc func(context.Context, *schematic.InstanceDiff, interface{}) error

func (r *Instance) create(ctx context.Context, d *InstanceData, meta interface{}) error {
	if r.Create != nil {
		if err := r.Create(d, meta); err != nil {
			return err
		}
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(TimeoutCreate))
	defer cancel()
	return nil
}

func (r *Instance) read(ctx context.Context, d *InstanceData, meta interface{}) error {
	if r.Read != nil {
		if err := r.Read(d, meta); err != nil {
			return err
		}
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(TimeoutRead))
	defer cancel()
	return nil
}

func (r *Instance) update(ctx context.Context, d *InstanceData, meta interface{}) error {
	if r.Update != nil {
		if err := r.Update(d, meta); err != nil {
			return err
		}
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(TimeoutUpdate))
	defer cancel()
	return nil
}

func (r *Instance) delete(ctx context.Context, d *InstanceData, meta interface{}) error {
	if r.Delete != nil {
		if err := r.Delete(d, meta); err != nil {
			return err
		}
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(TimeoutDelete))
	defer cancel()
	return nil
}

// Apply creates, updates, and/or deletes a resource.
func (r *Instance) Apply(ctx context.Context, s *schematic.InstanceState, d *schematic.InstanceDiff, meta interface{}) (*schematic.InstanceState, error) {
	// TODO: Implement this method
	return nil, nil
}

func (r *Instance) Data(s *schematic.InstanceState) *InstanceData {
	result, err := schemaMap(r.Schema).Data(s, nil)
	if err != nil {
		// At the time of writing, this isn't possible (Data never returns
		// non-nil errors). We panic to find this in the future if we have to.
		// I don't see a reason for Data to ever return an error.
		panic(err)
	}

	// load the Resource timeouts
	result.timeouts = r.Timeouts
	if result.timeouts == nil {
		result.timeouts = &InstanceTimeout{}
	}

	return result
}

// Noop is a convenience implementation of resource function which takes
// no action and returns no error.
func Noop(*InstanceData, interface{}) error {
	return nil
}
