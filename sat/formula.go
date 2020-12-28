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

// NumClause :
func (cnf CNF) NumClause() int {
	return len(cnf)
}

// NumVar :
func (cnf CNF) NumVar() int {
	numVar := 0
	for _, clause := range cnf {
		for _, literal := range clause {
			variable := abs(literal)
			if numVar < variable {
				numVar = variable
			}
		}
	}
	return numVar
}
