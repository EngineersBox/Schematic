package state

type ParsedState struct {
	Variables map[string]*Variable
	Instances map[string]*InstanceData
	Captures  map[string]*Capture
	Data      map[string]*Data
}
