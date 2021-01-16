package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/EngineersBox/ModularCLI/cli"
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

var commands = map[string]cli.SubCommand{
	"plan": {
		ErrorHandler: flag.ExitOnError,
		Arguments: []cli.Argument{
			{
				Type:         cli.TypeString,
				Name:         "sch",
				DefaultValue: "",
				HelpMsg:      "Schematic file (*.sch) to preview changes for",
				Required:     true,
				Optional:     false,
			},
			{
				Type:         cli.TypeBool,
				Name:         "dtf",
				DefaultValue: false,
				HelpMsg:      "Whether to output diff to file",
				Required:     false,
				Optional:     true,
			},
			{
				Type:         cli.TypeString,
				Name:         "diffout",
				DefaultValue: "diff_out",
				HelpMsg:      "File address to output diff to (requires --dtf=true) [default: diff_out]",
				Required:     false,
				Optional:     true,
			},
		},
	},
	"apply": {
		ErrorHandler: flag.ExitOnError,
		Arguments: []cli.Argument{
			{
				Type:         cli.TypeString,
				Name:         "sch",
				DefaultValue: "",
				HelpMsg:      "Schematic file (*.sch) to apply changes for",
				Required:     true,
				Optional:     false,
			},
			{
				Type:         cli.TypeBool,
				Name:         "diff",
				DefaultValue: true,
				HelpMsg:      "Whether to show diff for changes [default: true]",
				Required:     false,
				Optional:     true,
			},
		},
	},
}

func main() {
	schematicCli, err := cli.CreateCLI(commands)
	if err != nil {
		log.Fatal(err)
	}
	err = schematicCli.Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*schematicCli.Commands["plan"].Flags["sch"].GetString())
	fmt.Println(CapsuleConfig.Schema["containerId"])
}
