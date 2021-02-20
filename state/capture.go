package state

type Capture struct {
	Name          string
	HasDependency []string
	Handler       string
	Tainted       bool `json:"tainted"`
}
