package cli

type ArgType int

const (
	TypeInvalid ArgType = iota
	TypeInt
	TypeBool
	TypeString
)

type Argument struct {
	Type         ArgType
	Name         string
	DefaultValue interface{}
	Optional     bool
	Required     bool
	HelpMsg      string
	ValidateFunc ArgValidateFunc
}

type ArgValidateFunc func(interface{}) error

type TypedArgument interface {
	GetString() *string
	GetBool() *bool
	GetInt() *int
}

type StringArgument struct {
	Value *string
}

func (s StringArgument) GetString() *string {
	return s.Value
}

func (StringArgument) GetBool() *bool { return nil }
func (StringArgument) GetInt() *int   { return nil }

type IntArgument struct {
	Value *int
}

func (i IntArgument) GetInt() *int {
	return i.Value
}

func (IntArgument) GetBool() *bool     { return nil }
func (IntArgument) GetString() *string { return nil }

type BoolArgument struct {
	Value *bool
}

func (b BoolArgument) GetBool() *bool {
	return b.Value
}

func (BoolArgument) GetString() *string { return nil }
func (BoolArgument) GetInt() *int       { return nil }
