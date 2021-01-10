package main

import (
	"fmt"

	"github.com/EngineersBox/Schematic/schema"
)

func main() {
	containerSchema := map[string]*schema.Schema{
		"container_id": {
			Type:     schema.TypeString,
			Computed: false,
		},
	}
	fmt.Println(containerSchema["container_id"])
}
