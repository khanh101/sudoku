package sat

// Literal :
type Literal = int

// Clause :
type Clause []Literal

// CNF :
type CNF []Clause

// NewCNF :
func NewCNF() CNF {
	return make([]Clause, 0)
}

// Copy :
func (cnf CNF) Copy() CNF {
	out := make([]Clause, len(cnf))
	for ic, c := range cnf {
		out[ic] = make([]int, len(c))
		copy(out[ic], c)
	}
	return out
}
