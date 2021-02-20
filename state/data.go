package state

import (
	"strings"
)

type Data struct {
	Type       string
	Name       string
	Reference  string
	Attributes map[string]interface{}
	Tainted    bool `json:"tainted"`
}

func (d *Data) validateReference() bool {
	return strings.HasPrefix(d.Reference, "L::") || strings.HasPrefix(d.Reference, "W::")
}
