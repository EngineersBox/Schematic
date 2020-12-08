package schematic

import (
	"fmt"
	"strconv"
	"strings"
)

const VarReferencePrefix = "var."

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

// ---- REFERENCE ----

func resolveVariableReference(value string references map[string]Variable) (Variable, error) {
	if err := isVaraibleReferenceWellFormed(value, references); err != nil {
		return Variable{}, err
	}
	splitted []string := strings.Split(value, ".")
	ref := references[splitted[1]]
	if ref != nil {
		return *ref, nil
	}
	return Variable{}, &VariableError{fmt.Sprintf("No such variable exists: %s", splitted[1])}
}

func isVaraibleReferenceWellFormed(value string, references map[string]Variable) error {
	if len(value) < 1 || !isVariableReference(string) {
		return &VariableError{fmt.Sprintf("Not a variable reference: %s", value)}
	}
	splitted []string := strings.Split(value, ".")
	if len(splitted) != 2 {
		return &Variable{fmt.Sprintf("Invalid variable reference: %s", value)}
	}
}

func isVariableReference(value string) bool {
	return strings.HasPrefix(value, VarReferencePrefix)
}
