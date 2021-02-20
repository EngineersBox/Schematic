package parser

import (
	"fmt"
	"github.com/EngineersBox/Schematic/collection"
	"github.com/EngineersBox/Schematic/providers"
	"github.com/EngineersBox/Schematic/schema"
	"github.com/EngineersBox/Schematic/state"
	"strconv"
	"strings"
)

const (
	variableReferencePrefix = "var."
	fieldNestingDelimiter   = "->"
)

func (p *Parser) parseInstance(schem *state.ParsedState) (string, *state.InstanceState, error) {
	newInst := &state.InstanceState{}
	tok, lit := p.scanLine()
	if tok[0] != IDENT || !strings.Contains(lit[0], providerReferenceDelimiter) {
		return "", nil, fmt.Errorf("invalid provider reference for instance: " + lit[0])
	}
	provRef, err := newProviderReference(lit[0])
	if err != nil {
		return "", nil, err
	}
	if tok[1] != IDENT {
		return "", nil, fmt.Errorf("not a valid variable name")
	}
	newInst.ID = lit[1]
	if tok[2] != OPENBRACE {
		return "", nil, fmt.Errorf("missing open brace in instance declaration")
	}
	err = p.parseInstanceBody(newInst, provRef, schem)
	if err != nil {
		return "", nil, err
	}
	return lit[1], newInst, nil
}

func (p *Parser) parseInstanceBody(newInst *state.InstanceState, providerReference *ProviderReference, schem *state.ParsedState) error {
	provider := providers.InstalledProviders[providerReference.Provider]
	if provider == nil {
		return fmt.Errorf("not such provider: %s", providerReference.Provider)
	}
	instanceReference := provider.InstancesMap[providerReference.Kind]
	if instanceReference == nil {
		return fmt.Errorf("provider [%s] has no instance defintion: %s", providerReference.Provider, providerReference.Kind)
	}
	instanceSchema := instanceReference.Schema
	newInst.Provider = providerReference.Provider
	newInst.Type = providerReference.Kind
	newInst.Attributes = make(map[string]interface{})
	blockDepth := 0
	var nesting []string
	for {
		tokens, literals := p.scanLine()
		if len(tokens) == 1 && tokens[0] == CLOSEDBRACE {
			if blockDepth == 0 {
				break
			}
			blockDepth--
			continue
		}
		if len(tokens) != 3 {
			return fmt.Errorf("invalid assignment in instance body: \"%s\"", strings.Join(literals, " "))
		}
		if tokens[1] != EQUALS {
			return fmt.Errorf("assignment must be via equals operator. Invalid assignment: %s", strings.Join(literals, " "))
		}
		if tokens[2] == OPENBRACE {
			nesting = append(nesting, literals[0])
			blockDepth++
			continue
		}
		err := updateInstanceFields(nesting, literals, instanceSchema, schem, providerReference, newInst)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateInstanceFields(nesting []string, literals []string, instanceSchema map[string]*schema.Schema, schem *state.ParsedState, providerReference *ProviderReference, newInst *state.InstanceState) error {
	currentNesting := append(nesting, literals[0])
	if !validateSchemaFields(currentNesting, instanceSchema) {
		return fmt.Errorf(
			"instance [%s] has no schema field for: %s",
			providerReference.AsString(),
			strings.Join(currentNesting, fieldNestingDelimiter),
		)
	}
	assignableValue := literals[2]
	if strings.Contains(assignableValue, variableReferencePrefix) {
		variableReference := strings.TrimPrefix(assignableValue, variableReferencePrefix)
		variableType := schem.Variables[variableReference].BaseType
		if variableType == schematic.TypeInt {
			val, _ := schem.Variables[variableReference].Value.AsBigFloat().Int64()
			assignableValue = strconv.FormatInt(val, 10)
		} else if variableType == schematic.TypeFloat {
			val, _ := schem.Variables[variableReference].Value.AsBigFloat().Float64()
			assignableValue = fmt.Sprintf("%f", val)
		} else if variableType == schematic.TypeBool {
			assignableValue = schem.Variables[variableReference].Value.AsString()
		} else {
			return fmt.Errorf("unknown variable type: %v", variableType)
		}
	}
	updatedFields, err := recurseAssign(currentNesting, assignableValue, newInst.Attributes)
	if err != nil {
		return err
	}
	newInst.Attributes = updatedFields
	return nil
}
