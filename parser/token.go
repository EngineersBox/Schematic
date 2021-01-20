package parser

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT

	// Composites
	BLOCK
	ARRAY

	// Misc characters
	OPENBRACE     // {
	CLOSEDBRACE   // }
	OPENBRACKET   // [
	CLOSEDBRACKET // ]
	OPENPAREN     // (
	CLOSEDPAREN   // )
	COMMA         // ,
	PERIOD        // .

	// Keywords
	INSTANCE
	DATA
	CAPTURE
	VARIABLE
)
