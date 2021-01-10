package schematic

// EnumItem ... Enum entry pair
type EnumItem struct {
	index int
	value string
}

// Enum ... Enum type
type Enum struct {
	items []EnumItem
}

// Value ... Get the value of a given enum indexer
func (enum Enum) Value(findIndex int) string {
	for _, item := range enum.items {
		if item.index == findIndex {
			return item.value
		}
	}
	return "ID not found"
}

// Index ... Get the index from a given value
func (enum Enum) Index(findValue string) int {
	for idx, item := range enum.items {
		if findValue == item.value {
			return idx
		}
	}
	return -1
}

// Last ... Get the last (highest) index entry
func (enum Enum) Last() (int, string) {
	n := len(enum.items)
	return n - 1, enum.items[n-1].value
}
