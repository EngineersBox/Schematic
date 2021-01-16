package main

import (
	"fmt"

	"github.com/EngineersBox/Schematic/schema"
)

var CapsuleConfig = &schema.Instance{
	Create:      nil,
	Read:        nil,
	Update:      nil,
	Delete:      nil,
	Description: "A Capsule resource instance",
	Schema: map[string]*schema.Schema{
		"inbuilt": {
			Type:     schema.TypeBool,
			Computed: false,
			Required: true,
		},
		"containerId": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
		},
		"config": {
			Type:     schema.TypeMap,
			Computed: false,
			Required: true,
			Elem: &schema.Instance{
				Schema: map[string]*schema.Schema{
					"pidsMax": {
						Type:     schema.TypeInt,
						Computed: false,
						Required: false,
					},
					"memMax": {
						Type:     schema.TypeInt,
						Computed: false,
						Required: false,
					},
					"netClsId": {
						Type:     schema.TypeInt,
						Computed: false,
						Required: false,
					},
					"terminateOnClose": {
						Type:     schema.TypeBool,
						Computed: false,
						Required: false,
					},
				},
			},
		},
	},
}

func main() {
	fmt.Println(CapsuleConfig.Schema["containerId"])
}
