package cnf

// Clause :
type Clause []int

// Formula :
type Formula []Clause

// New :
func New() Formula {
	return make([]Clause, 0)
}

// AddClause :
func (cnf Formula) AddClause(clause []int) Formula {
	return append(cnf, clause)
}

// Copy :
func (cnf Formula) Copy() Formula {
	out := make([]Clause, len(cnf))
	for ic, c := range cnf {
		out[ic] = make([]int, len(c))
		copy(out[ic], c)
	}
	return out
}
