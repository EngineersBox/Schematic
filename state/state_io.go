package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type State struct {
	Filename    string
	ParsedState *ParsedState
	oldState    *JSONSchmState
	newState    *JSONSchmState
}

func NewJSONSchmState() *JSONSchmState {
	return &JSONSchmState{}
}

type JSONModuleState struct {
	Path      string                    `json:"path"`
	Instances map[string]*InstanceState `json:"instances"`
	Captures  map[string]*Capture       `json:"captures"`
	Data      map[string]*Data          `json:"data"`
}

type JSONSchmState struct {
	Version int               `json:"version"`
	Modules []JSONModuleState `json:"modules"`
}

var (
	operationalDirectory             = "operational"
	stateOut                         = fmt.Sprintf("%s/state.json", operationalDirectory)
	stateDirectoryMode   os.FileMode = 0700
	stateFileMode        os.FileMode = 0644
)

func createDirIfNotExists(path string, mode os.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, mode)
	}
	return nil
}

// State Format
// "JSONSchmState": {
//	  "version": "<INT>",
//	  "modules": [
//		  {
//			  "path": "<STRING>",
//			  "instances": [
//				  {
//					  "id": "<STRING>",
//					  "provider": "<STRING>",
//					  "name": "<STRING>",
//					  "type": "<STRING>",
//					  "attributes": {
//						  ..."<STRING>": "<STRING | INT | FLOAT | MAP>"
//					  },
//                    "meta": {
//						  ..."<STRING>": "<STRING | INT | FLOAT | MAP>"
//                    }
//				  }
//			  ],
//			  "captures": [
//				  {
//					  "name": "<STRING>",
//					  "hasDependency": [
//						  ..."<STRING>"
//					  ],
//					  "handler": "<STRING>"
//				  }
//			  ],
//			  "data": [
//				  "type": "<STRING>",
//				  "name": "<STRING>",
//				  "reference": "<L | W>::<STRING>"
//				  "schema": {
//					  ..."<STRING>": "<STRING | INT | FLOAT | MAP>"
//				  }
//			  ]
//		  }
//	  ]
// }

func (s *State) Init(parsedState *ParsedState) error {
	s.ParsedState = parsedState
	err := s.Read()
	if err != nil {
		return err
	}

	// TODO: Implement this

	return nil
}

func getStateFromInstanceData(instances map[string]*InstanceData) map[string]*InstanceState {
	states := make(map[string]*InstanceState)
	for k, v := range instances {
		if v.newState == nil {
			states[k] = v.state
			continue
		}
		states[k] = v.newState
	}
	return states
}

func (s *State) Write() error {
	// 1. Go through each instance/capture/data
	// 2. Update the JSONSchmState with the new instance/capture/data state
	// 3. Write the state to file

	s.newState = &JSONSchmState{
		Version: 1,
		Modules: make([]JSONModuleState, 0),
	}

	//if !strings.HasSuffix(schematicFilePath, ".schm") {
	//	givenFileType := strings.Split(schematicFilePath, ".")
	//	return fmt.Errorf(
	//		"schematic file path is not of a valid file type. Expected: .schm, Got: .%s",
	//		givenFileType[len(givenFileType) - 1],
	//	)
	//}

	s.newState.Modules = append(s.newState.Modules, JSONModuleState{
		Path:      "schematicFilePath",
		Instances: getStateFromInstanceData(s.ParsedState.Instances),
		Captures:  s.ParsedState.Captures,
		Data:      s.ParsedState.Data,
	})
	fmt.Println(s.newState.Modules[0].Instances["test_capsule"])
	file, _ := json.MarshalIndent(s.newState, "", "\t")
	err := createDirIfNotExists(operationalDirectory, stateDirectoryMode)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(s.Filename, file, stateFileMode)
	if err != nil {
		panic(err)
	}
	return nil
}

func (s *State) Read() error {
	// 1. Read the state from file
	// 2. Parse JSON to struct
	// 3. Check or errors in the state
	// 4. Generate a ParsedState
	// 5. Check existing resources for changes
	// 6. Mark resources as tainted if different

	data, err := ioutil.ReadFile(s.Filename)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, s.oldState)
	if err != nil {
		panic(err)
	}
	log.Println(s.oldState.Modules[0].Instances)
	return nil
}
