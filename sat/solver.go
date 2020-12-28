package sat

import (
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
)

// SolveOnce :
func SolveOnce(formula CNF) (satisfiable bool, assignment Assignment) {
	var2zVar := func(v int) z.Var {
		return z.Var(v)
	}
	numClause := formula.NumClause()
	numVar := formula.NumVar()

	g := gini.NewVc(numVar, numClause)
	for _, c := range formula {
		for _, v := range c {
			if v > 0 {
				g.Add(var2zVar(v).Pos())
			} else {
				g.Add(var2zVar(-v).Neg())
			}
		}
		g.Add(0)
	}
	satisfiable = (g.Solve() == 1)
	if satisfiable {
		assignment = NewAssignment(numVar)
		for v := range assignment {
			zVar := var2zVar(v)
			if g.Value(zVar.Pos()) {
				assignment[v] = +1
			} else {
				assignment[v] = -1
			}
		}
	}
	return satisfiable, assignment
}
