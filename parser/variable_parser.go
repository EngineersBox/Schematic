package parser

import (
	"fmt"
	"github.com/EngineersBox/Schematic/schema"
	"github.com/EngineersBox/Schematic/state"
	"github.com/zclconf/go-cty/cty"
	"strconv"
)

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
