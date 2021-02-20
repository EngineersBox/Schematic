package main

import (
	"flag"
	"fmt"
	"github.com/EngineersBox/ModularCLI/cli"
	"github.com/EngineersBox/Schematic/collection"
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
			Type:     schematic.TypeBool,
			Computed: false,
			Required: true,
		},
		"containerId": {
			Type:     schematic.TypeString,
			Computed: false,
			Required: true,
		},
		"config": {
			Type:     schematic.TypeMap,
			Computed: false,
			Required: true,
			Elem: map[string]*schema.Schema{
				"pidsMax": {
					Type:     schematic.TypeInt,
					Computed: false,
					Required: false,
				},
				"memMax": {
					Type:     schematic.TypeInt,
					Computed: false,
					Required: false,
				},
				"netClsId": {
					Type:     schematic.TypeInt,
					Computed: false,
					Required: false,
				},
				"terminateOnClose": {
					Type:     schematic.TypeBool,
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
		Flags: []*cli.Flag{
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
				Name:     "schm",
				Position: 0,
				ValidateFunc: func(arg cli.Parameter) error {
					if !strings.Contains(*arg.GetString(), ".schm") {
						return fmt.Errorf("filetype must be .sch")
					}
					return nil
				},
			},
		},
	},
	"apply": {
		ErrorHandler: flag.ExitOnError,
		Flags: []*cli.Flag{
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
				Name:     "schm",
				Position: 0,
				ValidateFunc: func(arg cli.Parameter) error {
					if !strings.Contains(*arg.GetString(), ".schm") {
						return fmt.Errorf("filetype must be .schm")
					}
					return nil
				},
			},
		},
	},
	"install": {
		ErrorHandler: flag.ExitOnError,
		Flags:        nil,
		Parameters: []*cli.Parameter{
			{
				Type:     cli.TypeString,
				Name:     "install_provider",
				Position: 0,
				ValidateFunc: func(arg cli.Parameter) error {
					_, exists := internalRegistry[*arg.GetString()]
					if !exists {
						return fmt.Errorf("invalid provider: %s", *arg.GetString())
					}
					return nil
				},
			},
		},
	},
}

var (
	schemaReferences     = make(map[string]interface{})
	operationalDirectory = "operational"
	stateOut             = fmt.Sprintf("%s/state.json", operationalDirectory)
	internalRegistry     = map[string]string{
		"capsule": "github.com/engineersbox/terraform_provider_capsule",
		"aws":     "github.com/aws/terraform_provider_aws",
	}
)

func createDirIfNotExists(path string, mode os.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, mode)
	}
	return nil
}

func main() {
	providers.InstalledProviders["capsule"] = &providers.Provider{}
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

	f, err := os.Open(*schematicCli.Commands["plan"].Params["schm"].GetString())
	if err != nil {
		panic(err)
	}

	p := parser.NewParser(f)
	ps, err := p.Parse()
	if err != nil {
		panic(err)
	}

	newState := &state.State{
		Filename: stateOut,
	}
	err = newState.Init(ps)
	if err != nil {
		log.Fatal(err)
	}

	newVar := ps.Variables["testcapsule_clsid"]
	if newVar.BaseType == schematic.TypeFloat || newVar.BaseType == schematic.TypeInt {
		fmt.Println(newVar.Value.AsBigFloat())
	} else if newVar.BaseType == schematic.TypeString {
		fmt.Println(newVar.Value.AsString())
	}

	value, err := ps.Instances["test_capsule"].GetFromNesting([]string{"config", "netClsId"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("test_capsule->config->netClsId: %s\n", value)

	value, err = ps.Instances["test_capsule"].GetFromNesting([]string{"containerId"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("test_capsule->containerId: %s\n", value)

	//file, _ := json.MarshalIndent(ps, "", " ")
	//err = createDirIfNotExists(operationalDirectory, 0700)
	//if err != nil {
	//	panic(err)
	//}
	//_ = ioutil.WriteFile(stateOut, file, 0644)
}
