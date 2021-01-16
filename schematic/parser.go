package schematic

import (
	"fmt"
)

var reservedKeywords = []string{
	"instance",
	"data",
	"capture",
	"structure",
	"variable",
}

// ---- ALL TYPES ----

type SchemaType string

var allTypes = []SchemaType{
	"String",
	"Integer",
	"Float",
	"Boolean",
	"Block",
	"Array",
}

// ---- BASIC TYPES ----

var BasicTypes = allTypes[:4]

// ---- COMPLEX TYPES ----

var ComplexTypes = allTypes[4:]

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
