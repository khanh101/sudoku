package sat

// Implication :
func Implication(formulaIn CNF, assignmentIn map[int]bool) (unsatisfiable bool, formulaOut CNF, assignmentOut map[int]bool) {
	assignmentOut = assignmentIn
	formulaOut = formulaIn
	abs := func(x int) int {
		if x > 0 {
			return +x
		}
		if x < 0 {
			return -x
		}
		return 0
	}
	for {
		better := false
		// assign to all truth value
		for _, clause := range formulaOut {
			if len(clause) == 1 {
				l := clause[0]
				vi := abs(l)
				if _, ok := assignmentOut[vi]; !ok {
					assignmentOut[vi] = l > 0
					better = true
				}
			}
		}
		// detect unsat and remove 1-clause
		oneClauseIdx := make([]int, 0)
		for cidx, clause := range formulaOut {
			unsatisfiable = true
			zeroLiteralIdx := make([]int, 0)
			for lidx, l := range clause {
				vi := abs(l)
				val, ok := assignmentOut[vi]
				// at least 1 unknown -> pass
				if !ok {
					unsatisfiable = false
					break
				}
				// true -> remove from formulaOut
				if l > 0 && val == true {
					unsatisfiable = false
					oneClauseIdx = append(oneClauseIdx, cidx)
					break
				}
				// true -> remove from formulaOut
				if l < 0 && val == false {
					unsatisfiable = false
					oneClauseIdx = append(oneClauseIdx, cidx)
					break
				}
				// false -> remove from claus
				zeroLiteralIdx = append(zeroLiteralIdx, lidx)
			}

			for i := len(zeroLiteralIdx) - 1; i >= 0; i-- {
				lidx := zeroLiteralIdx[i]
				formulaOut[cidx][lidx] = formulaOut[cidx][len(formulaOut[cidx])-1]
				formulaOut[cidx] = formulaOut[cidx][:len(formulaOut[cidx])-1]
				better = true
			}
			if unsatisfiable {
				return unsatisfiable, formulaOut, assignmentOut
			}
		}
		for i := len(oneClauseIdx) - 1; i >= 0; i-- {
			cidx := oneClauseIdx[i]
			formulaOut[cidx] = formulaOut[len(formulaOut)-1]
			formulaOut = formulaOut[:len(formulaOut)-1]
			better = true
		}
		if !better {
			break
		}
	}
	return unsatisfiable, formulaOut, assignmentOut
}
