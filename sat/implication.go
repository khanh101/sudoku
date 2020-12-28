package sat

// Explanation :
type Explanation map[Literal]int

// Implication :
func Implication(formula CNF, bootstrap Assignment, explain bool) (unsatisfiable bool, assignment Assignment, explanation Explanation) {
	if explain {
		explanation = make(Explanation)
	}

	numVar := formula.NumVar()
	assignment = NewAssignment(numVar)
	for i, v := range bootstrap {
		assignment[i] = v
	}
	literalValue := func(literal Literal) Value {
		return assignment[abs(literal)] * sign(literal)
	}
	clauseValue := func(clause Clause) (Value, int, int) {
		// 1 if clause is sat, -1 if clause is unsat, (0, numZero, firstZeroIdx) if there still chance
		numZero := 0
		firstZeroIdx := -1
		for idx, literal := range clause {
			v := literalValue(literal)
			if v == ValueTrue {
				return ValueTrue, 0, 0
			}
			if v == ValueUnknown {
				numZero++
				if firstZeroIdx == -1 {
					firstZeroIdx = idx
				}
			}
		}
		if numZero > 0 {
			return ValueUnknown, numZero, firstZeroIdx
		}
		return ValueFalse, 0, 0
	}

	for {
		unitprop := false
		for idx, clause := range formula {
			value, numZero, firstZeroIdx := clauseValue(clause)
			if value == ValueFalse {
				unsatisfiable = true
				return unsatisfiable, assignment, explanation
			}
			if value == ValueTrue {
				continue
			}
			if value == ValueUnknown && numZero == 1 {
				activateLiteral := clause[firstZeroIdx]
				assignment[abs(activateLiteral)] = sign(activateLiteral)
				unitprop = true
				if explain {
					explanation[activateLiteral] = idx
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
