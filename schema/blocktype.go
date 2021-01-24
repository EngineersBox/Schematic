package schema

type BlockType int

const (
	InvalidBlockType BlockType = iota
	TypeVariable
	TypeInstance
	TypeData
	TypeCapture
)

func (b *BlockType) ToString() string {
	switch *b {
	case TypeVariable:
		return "variable"
	case TypeInstance:
		return "instance"
	case TypeData:
		return "data"
	case TypeCapture:
		return "capture"
	case InvalidBlockType:
		return "invalid"
	default:
		return "invalid"
	}
}

func ToBlockType(s string) BlockType {
	if s == "variable" {
		return TypeVariable
	} else if s == "instance" {
		return TypeInstance
	} else if s == "data" {
		return TypeData
	} else if s == "capture" {
		return TypeCapture
	} else {
		return InvalidBlockType
	}
}
