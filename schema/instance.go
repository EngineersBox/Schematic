package schema

import (
	"context"
	"github.com/EngineersBox/Schematic/state"
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

type CreateFunc func(*state.InstanceData, interface{}) error
type ReadFunc func(*state.InstanceData, interface{}) error
type UpdateFunc func(*state.InstanceData, interface{}) error
type DeleteFunc func(*state.InstanceData, interface{}) error
type ExistsFunc func(*state.InstanceData, interface{}) (bool, error)

type CustomizeDiffFunc func(context.Context, *state.InstanceDiff, interface{}) error

func (r *Instance) create(ctx context.Context, d *state.InstanceData, meta interface{}) error {
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

func (r *Instance) read(ctx context.Context, d *state.InstanceData, meta interface{}) error {
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

func (r *Instance) update(ctx context.Context, d *state.InstanceData, meta interface{}) error {
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

func (r *Instance) delete(ctx context.Context, d *state.InstanceData, meta interface{}) error {
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
func (r *Instance) Apply(ctx context.Context, s *Instance, d *state.InstanceDiff, meta interface{}) error {
	// TODO: Implement this method
	return nil
}

func (r *Instance) Data(s *InstanceState) *state.InstanceData {
	// TODO: Implement this
	return nil
}

// Noop is a convenience implementation of resource function which takes
// no action and returns no error.
func Noop(*state.InstanceData, interface{}) error {
	return nil
}
