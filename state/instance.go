package state

type Instance struct {
	Provider string                 `json:"provider"`
	Type     string                 `json:"type"`
	Name     string                 `json:"name"`
	Fields   map[string]interface{} `json:"fields"`
}
