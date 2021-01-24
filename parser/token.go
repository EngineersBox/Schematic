package parser

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT

	// Misc characters
	OPENBRACE     // {
	CLOSEDBRACE   // }
	OPENBRACKET   // [
	CLOSEDBRACKET // ]
	OPENPAREN     // (
	CLOSEDPAREN   // )
	COMMA         // ,
	PERIOD        // .
	EQUALS        // =

	// Keywords
	INSTANCE
	DATA
	CAPTURE
	VARIABLE
)
