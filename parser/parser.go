package parser

import (
	"github.com/EngineersBox/Schematic/state"
	"io"
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

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan(returnOnNL bool) (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan(returnOnNL)

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
func (p *Parser) scanIgnoreWhitespace(returnOnNL bool) (tok Token, lit string) {
	for {
		tok, lit = p.scan(returnOnNL)
		if (returnOnNL && strings.ContainsRune(lit, '\n')) || tok != WS {
			break
		}
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
