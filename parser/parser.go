package parser

import (
	"fmt"
	"github.com/EngineersBox/Schematic/providers"
	"github.com/EngineersBox/Schematic/schema"
	"github.com/EngineersBox/Schematic/state"
	"github.com/zclconf/go-cty/cty"
	"io"
	"strconv"
	"strings"
)

type ParsedState struct {
	Variables map[string]*state.Variable
	Instances map[string]*state.Instance
}

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

const (
	variableReferencePrefix = "var."
	fieldNestingDelimiter   = "->"
)

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses tokens into a map of declarations
func (p *Parser) Parse() (*ParsedState, error) {
	schem := &ParsedState{
		Variables: make(map[string]*state.Variable),
		Instances: make(map[string]*state.Instance),
	}
	for {
		token, _ := p.scanIgnoreWhitespace(false)
		if token == EOF {
			break
		}
		switch token {
		case VARIABLE:
			name, newVar, err := p.parseVariable()
			if err != nil {
				return nil, err
			}
			schem.Variables[name] = newVar
		case INSTANCE:
			name, newInstance, err := p.parseInstance(schem)
			if err != nil {
				return nil, err
			}
			schem.Instances[name] = newInstance
		}
	}
	return schem, nil
}

func (p *Parser) parseVariable() (string, *state.Variable, error) {
	newVar := &state.Variable{}
	token, name := p.scanIgnoreWhitespace(false)
	if token != IDENT {
		return "", nil, fmt.Errorf("not a valid variable name")
	}
	newVar.Name = name
	tok, lit := p.scanIgnoreWhitespace(false)
	if tok != OPENBRACE {
		return "", nil, fmt.Errorf("missing open brace in variable declaration")
	}
	tok, lit = p.scanIgnoreWhitespace(false)
	if tok != IDENT || lit != "value" {
		return "", nil, fmt.Errorf("variable body must only have value declaration")
	}
	tok, lit = p.scanIgnoreWhitespace(false)
	if tok != EQUALS {
		return "", nil, fmt.Errorf("missing assignment operator '=' in variable declaration")
	}
	tok, lit = p.scanIgnoreWhitespace(false)
	if tok != IDENT {
		return "", nil, fmt.Errorf("variable value must only contain a literal")
	}
	val, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		val, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			newVar.Value = cty.StringVal(lit)
			newVar.BaseType = schema.TypeString
		} else {
			newVar.Value = cty.NumberFloatVal(val)
			newVar.BaseType = schema.TypeFloat
		}
	} else {
		newVar.Value = cty.NumberIntVal(val)
		newVar.BaseType = schema.TypeInt
	}
	tok, lit = p.scanIgnoreWhitespace(false)
	if tok != CLOSEDBRACE {
		return "", nil, fmt.Errorf("missing closing brace in variable declaration")
	}
	return name, newVar, nil
}

func (p *Parser) parseInstance(schem *ParsedState) (string, *state.Instance, error) {
	newInst := &state.Instance{}
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
	newInst.Name = lit[1]
	if tok[2] != OPENBRACE {
		return "", nil, fmt.Errorf("missing open brace in instance declaration")
	}
	err = p.parseInstanceBody(newInst, provRef, schem)
	if err != nil {
		return "", nil, err
	}
	return lit[1], newInst, nil
}

func (p *Parser) parseInstanceBody(newInst *state.Instance, providerReference *ProviderReference, schem *ParsedState) error {
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
	newInst.Fields = make(map[string]interface{})
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
		if len(nesting) == 0 {
			schemaField := instanceSchema[literals[0]]
			if schemaField == nil {
				return fmt.Errorf(
					"no such field [%s] for instance [%s]",
					literals[0],
					providerReference.AsString(),
				)
			}
			assignableValue := literals[2]
			if strings.Contains(assignableValue, variableReferencePrefix) {
				variableReference := strings.TrimPrefix(assignableValue, variableReferencePrefix)
				assignableValue = schem.Variables[variableReference].Value.AsString()
			}
			currentNesting := append(nesting, literals[0])
			if !validateSchemaFields(currentNesting, instanceSchema) {
				return fmt.Errorf(
					"instance [%s] has no schema field for: %s",
					providerReference.AsString(),
					strings.Join(currentNesting, fieldNestingDelimiter),
				)
			}
			newInst.Fields[literals[0]] = assignableValue
		} else {
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
				if variableType == schema.TypeInt {
					val, _ := schem.Variables[variableReference].Value.AsBigFloat().Int64()
					assignableValue = strconv.FormatInt(val, 10)
				} else if variableType == schema.TypeFloat {
					val, _ := schem.Variables[variableReference].Value.AsBigFloat().Float64()
					assignableValue = fmt.Sprintf("%f", val)
				} else if variableType == schema.TypeBool {
					assignableValue = schem.Variables[variableReference].Value.AsString()
				} else {
					return fmt.Errorf("unknown variable type: %v", variableType)
				}
			}
			updatedFields, err := recurseAssign(currentNesting, assignableValue, newInst.Fields)
			if err != nil {
				return err
			}
			newInst.Fields = updatedFields
		}
	}
	return nil
}

func validateSchemaFields(fields []string, instanceSchema map[string]*schema.Schema) bool {
	hasField := instanceSchema[fields[0]] != nil
	if len(fields) == 1 {
		return hasField
	}
	if instanceSchema[fields[0]].Type != schema.TypeMap {
		return false
	}
	nestedSchema, ok := instanceSchema[fields[0]].Elem.(map[string]*schema.Schema)
	if !ok {
		return false
	}
	return validateSchemaFields(fields[1:], nestedSchema)
}

func recurseAssign(nesting []string, value interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	if len(nesting) == 1 {
		obj[nesting[0]] = value
		return obj, nil
	}
	next := nesting[0]
	if obj[next] == nil {
		obj[next] = make(map[string]interface{})
	}
	nextMap, ok := obj[next].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not update nesting")
	}
	recusedValue, err := recurseAssign(nesting[1:], value, nextMap)
	obj[next] = recusedValue
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan(returnOnCR bool) (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan(returnOnCR)

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

func (p *Parser) scanLine() (tokens []Token, literals []string) {
	for {
		tok, lit := p.scan(true)
		if tok == EOF {
			break
		}
		if strings.ContainsRune(lit, '\n') {
			break
		}
		if tok != WS {
			tokens = append(tokens, tok)
			literals = append(literals, lit)
		}
	}
	return
}

// scanIgnoreWhitespace scans the next non-whitespace tokens.
func (p *Parser) scanIgnoreWhitespace(returnOnCR bool) (tok Token, lit string) {
	for {
		tok, lit = p.scan(returnOnCR)
		if (returnOnCR && strings.ContainsRune(lit, '\n')) || tok != WS {
			break
		}
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
