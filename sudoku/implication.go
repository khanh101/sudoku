package sudoku

import "github.com/khanhhhh/sudoku/sat"

func sign(x int) int {
	if x > 0 {
		return +1
	}
	if x < 0 {
		return -1
	}
	return 0
}

// Implication : implication but stop in first positive literal with len of clause > 1
func Implication(formula sat.CNF, explain bool) (unsatisfiable bool, assignment sat.Assignment, explanation sat.Explanation) {
	if explain {
		explanation = make(sat.Explanation)
	}

	numVar := formula.NumVar()
	assignment = sat.NewAssignment(numVar)
	literalValue := func(literal sat.Literal) sat.Value {
		return assignment[abs(literal)] * sign(literal)
	}
	clauseValue := func(clause sat.Clause) (sat.Value, int, int) {
		// 1 if clause is sat, -1 if clause is unsat, (0, numZero, firstZeroIdx) if there still chance
		numZero := 0
		firstZeroIdx := -1
		for idx, literal := range clause {
			v := literalValue(literal)
			if v == sat.ValueTrue {
				return sat.ValueTrue, 0, 0
			}
			if v == sat.ValueUnknown {
				numZero++
				if firstZeroIdx == -1 {
					firstZeroIdx = idx
				}
			}
		}
		if numZero > 0 {
			return sat.ValueUnknown, numZero, firstZeroIdx
		}
		return sat.ValueFalse, 0, 0
	}

	for {
		unitprop := false
		for idx, clause := range formula {
			value, numZero, firstZeroIdx := clauseValue(clause)
			if value == sat.ValueFalse {
				unsatisfiable = true
				return unsatisfiable, assignment, explanation
			}
			if value == sat.ValueTrue {
				continue
			}
			if value == sat.ValueUnknown && numZero == 1 {
				activateLiteral := clause[firstZeroIdx]
				assignment[abs(activateLiteral)] = sign(activateLiteral)
				unitprop = true
				if explain {
					explanation[activateLiteral] = idx
				}
				if activateLiteral > 0 && len(clause) > 1 {
					unsatisfiable = false
					return unsatisfiable, assignment, explanation
				}
			}
		}
		if !unitprop {
			break
		}
	}
	unsatisfiable = false
	return unsatisfiable, assignment, explanation
}
