package sat

// Clause :
type Clause []int

// CNF :
type CNF []Clause

// NewCNF :
func NewCNF() CNF {
	return make([]Clause, 0)
}

// AddClause :
func (cnf CNF) AddClause(clause []int) CNF {
	return append(cnf, clause)
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
