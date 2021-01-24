package parser

import (
	"fmt"
	"github.com/EngineersBox/Schematic/schema"
	"github.com/EngineersBox/Schematic/state"
	"github.com/zclconf/go-cty/cty"
	"io"
	"strconv"
	"strings"
)

const providerReferenceDelimiter = "::"

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
	token, _ := p.scanIgnoreWhitespace()
	switch token {
	case VARIABLE:
		name, newVar, err := p.parseVariable()
		if err != nil {
			return nil, err
		}
		schem.Variables[name] = newVar
	case INSTANCE:
		name, newInstance, err := p.parseInstance()
		if err != nil {
			return nil, err
		}
		schem.Instances[name] = newInstance
	}
	return schem, nil
}

func (p *Parser) parseVariable() (string, *state.Variable, error) {
	newVar := &state.Variable{}
	token, name := p.scanIgnoreWhitespace()
	if token != IDENT {
		return "", nil, fmt.Errorf("not a valid variable name")
	}
	newVar.Name = name
	tok, lit := p.scanIgnoreWhitespace()
	if tok != OPENBRACE {
		return "", nil, fmt.Errorf("missing open brace in variable declaration")
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT || lit != "value" {
		return "", nil, fmt.Errorf("variable body must only have value declaration")
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != EQUALS {
		return "", nil, fmt.Errorf("missing assignment operator '=' in variable declaration")
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return "", nil, fmt.Errorf("variable value must only contain a literal")
	}
	val, err := strconv.ParseFloat(lit, 10)
	if err != nil {
		val, err := strconv.ParseInt(lit, 10, 10)
		if err != nil {
			newVar.Value = cty.StringVal(lit)
			newVar.BaseType = schema.TypeString
		} else {
			newVar.Value = cty.NumberIntVal(val)
			newVar.BaseType = schema.TypeInt
		}
	} else {
		newVar.Value = cty.NumberFloatVal(val)
		newVar.BaseType = schema.TypeFloat
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != CLOSEDBRACE {
		return "", nil, fmt.Errorf("missing closing brace in variable declaration")
	}
	return name, newVar, nil
}

func (p *Parser) parseInstance() (string, *state.Instance, error) {
	newInst := &state.Instance{}
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT || strings.Contains(lit, providerReferenceDelimiter) {
		return "", nil, fmt.Errorf("invalid provider reference for instance")
	}
	provRef := strings.SplitAfter(lit, providerReferenceDelimiter)
	if len(provRef) != 2 {
		return "", nil, fmt.Errorf("provider reference must be of the form \"<provider>::<type>\"")
	}
	newInst.Provider = provRef[0]
	newInst.Type = provRef[1]
	tok, name := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return "", nil, fmt.Errorf("not a valid variable name")
	}
	newInst.Name = name
	tok, lit = p.scanIgnoreWhitespace()
	if tok != OPENBRACE {
		return "", nil, fmt.Errorf("missing open brace in instance declaration")
	}
	// TODO: Parse instance body
	tok, lit = p.scanIgnoreWhitespace()
	if tok != CLOSEDBRACE {
		return "", nil, fmt.Errorf("missing losing brace in instance declaration")
	}
	return name, newInst, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
