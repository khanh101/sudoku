package sudokusolver

import "github.com/khanhhhh/sudoku/satsolver"

// SolveOnce :
func SolveOnce(n Size, board Board, excludedList []Board) (sat bool, result Board) {
	formula := Reduce(n, board, excludedList)
	sat, assignment := satsolver.SolveOnce(formula)
	if sat {
		result = NewBoard(n)
		for vi, boolean := range assignment {
			if boolean {
				pi := v2p[n][vi]
				result[pi.row][pi.col] = pi.val
			}
		}
	}
	return sat, result
}
