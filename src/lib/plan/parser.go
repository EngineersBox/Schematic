package schematic

import (
	"fmt"
	"strconv"
)

var reservedKeywords []string = []string{
	"instance",
	"data",
	"capture",
	"structure",
	"variable",
}

// ---- ALL TYPES ----

type SchemaType string

var allTypes []SchemaType = []SchemaType{
	"String",
	"Integer",
	"Float",
	"Boolean",
	"Block",
	"Array",
}

// ---- BASIC TYPES ----

var BasicTypes []SchemaType = allTypes[:4]

// ---- VARIABLE ----

type Variable struct {
	value     string
	valueType SchemaType
}

type VariableError struct {
	prob string
}

func (e *VariableError) Error() string {
	return e.prob
}

func (v *Variable) ValidateType() error {
	for _, t := range BasicTypes {
		if t == v.valueType {
			return nil
		}
	}
	return &VariableError{fmt.Sprintf("Invalid type: %s", v.valueType)}
}

func (v *Variable) ToBool() (bool, error) {
	if v.valueType != BasicTypes[4] {
		return false, &VariableError{fmt.Sprintf("Not of type 'Boolean': %s", v.valueType)}
	}
	val, err := strconv.ParseBool(v.value)
	if err != nil {
		return false, &VariableError{err.Error()}
	}
	return val, nil
}

func (v *Variable) ToInt() (int, error) {
	if v.valueType != BasicTypes[2] {
		return 0, &VariableError{fmt.Sprintf("Not of type 'Integer': %s", v.valueType)}
	}
	val, err := strconv.Atoi(v.value)
	if err != nil {
		return 0, &VariableError{err.Error()}
	}
	return val, nil
}

func (v *Variable) ToFloat64() (float64, error) {
	if v.valueType != BasicTypes[3] || len(v.value) < 1 || string(v.value[len(v.value)-1]) != "f" {
		return 0.0, &VariableError{fmt.Sprintf("Not of type 'Float': %s", v.valueType)}
	}
	val, err := strconv.ParseFloat(v.value[:len(v.value)-1], 64)
	if err != nil {
		return 0, &VariableError{err.Error()}
	}
	return val, nil
}

func (v *Variable) ToString() (string, error) {
	if v.valueType != BasicTypes[0] {
		return "", &VariableError{fmt.Sprintf("Not of type 'String': %s", v.valueType)}
	}
	return v.value, nil
}

// ---- COMPLEX TYPES ----

var ComplexTypes []SchemaType = allTypes[4:]

// ---- FIELD PAIR ----

type FieldPair struct {
	key       string
	value     string
	valueType SchemaType
}

// ---- BLOCK ----

type Block struct {
	Name  string
	Pairs []FieldPair
}

type BlockError struct {
	prob string
}

func (e *BlockError) Error() string {
	return e.prob
}

func (b *Block) Size() int {
	return len(b.Pairs)
}

func (b *Block) Get(keyStr string) (string, SchemaType, error) {
	for _, pair := range b.Pairs {
		if pair.key == keyStr {
			return pair.value, pair.valueType, nil
		}
	}
	return "", "", &BlockError{fmt.Sprintf("No field with key: %s", keyStr)}
}

// ---- ARRAY ENTRY ----

type ArrayEntry struct {
	value     string
	valueType SchemaType
}

// ---- ARRAY ----

type Array struct {
	Name   string
	Values []ArrayEntry
}

type ArrayError struct {
	prob string
}

func (e *ArrayError) Error() string {
	return e.prob
}

func (a *Array) Size() int {
	return len(a.Values)
}

func (a *Array) AtIdx(idx int) (string, SchemaType, error) {
	if a.Size() < 1 {
		return "", "", &ArrayError{"No elements in array"}
	}
	if idx > a.Size() {
		return "", "", &ArrayError{fmt.Sprintf("Index %d out of bounds for array length %d", idx, a.Size())}
	}
	return a.Values[idx].value, a.Values[idx].valueType, nil
}
