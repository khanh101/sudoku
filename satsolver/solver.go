package satsolver

import (
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
	"github.com/khanhhhh/sudoku/cnf"
)

func abs(value int) int {
	if value > 0 {
		return value
	}
	return -value
}

// SolveOnce :
func SolveOnce(cnf cnf.Formula) (sat bool, assignment map[int]bool) {
	numClause := len(cnf)
	var2zVar := make(map[int]z.Var)
	counter := z.Var(0)
	for _, c := range cnf {
		for _, v := range c {
			_, ok := var2zVar[abs(v)]
			if !ok {
				counter++
				var2zVar[abs(v)] = counter
			}
		}
	}
	numVar := int(counter)

	g := gini.NewVc(numVar, numClause)
	for _, c := range cnf {
		for _, v := range c {
			if v > 0 {
				g.Add(var2zVar[v].Pos())
			} else {
				g.Add(var2zVar[-v].Neg())
			}
		}
		g.Add(0)
	}
	sat = (g.Solve() == 1)
	if sat {
		assignment = make(map[int]bool)
		for v, zVar := range var2zVar {
			assignment[v] = g.Value(zVar.Pos())
		}
	}
	return sat, assignment
}
