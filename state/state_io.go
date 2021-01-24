package state

import (
	"fmt"
	"github.com/EngineersBox/Schematic/schema"
	"os"
)

type State struct {
	Filename string
	Current  map[string]map[string]interface{}
	Writer   Writer
	Reader   Reader
}

type Writer func(filename string, data []byte, perm os.FileMode) error
type Reader func(filename string) ([]byte, error)

func (s *State) GetByType(blockType schema.BlockType) (map[string]interface{}, error) {
	if blockType == schema.InvalidBlockType {
		return nil, fmt.Errorf("invalid block type when retrieving from state")
	}
	return s.Current[blockType.ToString()], nil
}
