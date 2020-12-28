package sat

// Value :
type Value = int

// Value Constants
const (
	ValueTrue    = +1
	ValueUnknown = 0
	ValueFalse   = -1
)

// Assignment :
type Assignment []Value

// NewAssignment :
func NewAssignment(numVar int) Assignment {
	return make([]Value, numVar+1)
}

// NumVar :
func (a Assignment) NumVar() int {
	return len(a) - 1
}
