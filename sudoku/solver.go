package sudoku

import "github.com/khanh-nguyen-code/sudoku/sat"

// SolveOnce :
func SolveOnce(n int, board Board, excludedList []Board) (satisfiable bool, result Board) {
	formula := Reduce(n, board, excludedList)
	satisfiable, assignment := sat.SolveOnce(formula)
	if satisfiable {
		result = NewBoard(n)
		for vi, value := range assignment {
			if value == sat.ValueTrue {
				pi := v2p[n][vi]
				result[pi.row][pi.col] = pi.val
			}
		}
	}
	return satisfiable, result
}

// SolveAll :
func SolveAll(n int, board Board) (res Result) {
	return &result{
		n:            n,
		board:        board,
		excludedList: make([]Board, 0),
	}
}
