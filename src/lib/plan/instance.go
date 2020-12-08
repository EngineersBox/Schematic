package schematic

type Instance struct {
	name          string
	instType      string
	hasDependency Array
	count         int
	structure     Structure
}
