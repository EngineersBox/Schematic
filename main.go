package main

import (
	"flag"
	"fmt"
	"github.com/EngineersBox/ModularCLI/cli"
	"github.com/EngineersBox/Schematic/parser"
	"github.com/EngineersBox/Schematic/providers"
	"github.com/EngineersBox/Schematic/schema"
	"github.com/EngineersBox/Schematic/state"
	"log"
	"os"
	"strings"
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
			Elem: map[string]*schema.Schema{
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
}

var commands = map[string]cli.SubCommand{
	"plan": {
		ErrorHandler: flag.ExitOnError,
		Arguments: []*cli.Argument{
			{
				Type:         cli.TypeBool,
				Name:         "dtf",
				DefaultValue: false,
				HelpMsg:      "Whether to output diff to file",
				Required:     false,
			},
			{
				Type:         cli.TypeString,
				Name:         "diffout",
				DefaultValue: "diff_out",
				HelpMsg:      "File address to output diff to (requires --dtf=true) [default: diff_out]",
				Required:     false,
			},
		},
		Parameters: []*cli.Parameter{
			{
				Type:     cli.TypeString,
				Name:     "sch",
				Position: 0,
				ValidateFunc: func(arg cli.Parameter) error {
					if !strings.Contains(*arg.GetString(), ".sch") {
						return fmt.Errorf("filetype must be .sch")
					}
					return nil
				},
			},
		},
	},
	"apply": {
		ErrorHandler: flag.ExitOnError,
		Arguments: []*cli.Argument{
			{
				Type:         cli.TypeBool,
				Name:         "diff",
				DefaultValue: true,
				HelpMsg:      "Whether to show diff for changes [default: true]",
				Required:     false,
			},
		},
		Parameters: []*cli.Parameter{
			{
				Type:     cli.TypeString,
				Name:     "sch",
				Position: 0,
				ValidateFunc: func(arg cli.Parameter) error {
					if !strings.Contains(*arg.GetString(), ".sch") {
						return fmt.Errorf("filetype must be .sch")
					}
					return nil
				},
			},
		},
	},
}

var schemaReferences = make(map[string]interface{})

func main() {
	providers.InstalledProviders["capsule"] = &schema.Provider{}
	providers.InstalledProviders["capsule"].InstancesMap = make(map[string]*schema.Instance)
	providers.InstalledProviders["capsule"].InstancesMap["config"] = CapsuleConfig

	schemaReferences["capsule::config"] = *CapsuleConfig

	schematicCli, err := cli.CreateCLI(commands)
	if err != nil {
		log.Fatal(err)
	}
	err = schematicCli.Parse()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(*schematicCli.Commands["plan"].Params["sch"].GetString())
	if err != nil {
		panic(err)
	}

	p := parser.NewParser(f)
	ps, err := p.Parse()
	if err != nil {
		panic(err)
	}

	newVar := ps.Variables["testcapsule_clsid"]
	if newVar.BaseType == schema.TypeFloat || newVar.BaseType == schema.TypeInt {
		fmt.Println(newVar.Value.AsBigFloat())
	} else if newVar.BaseType == schema.TypeString {
		fmt.Println(newVar.Value.AsString())
	}
	value, err := state.GetInstanceField([]string{"config", "netClsId"}, ps.Instances["test_capsule"].Fields)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("test_capsule->config->netClsId: %s\n", value)
}
