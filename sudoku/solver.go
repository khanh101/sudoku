package sudoku

import "github.com/khanhhhh/sudoku/sat"

// SolveOnce :
func SolveOnce(n Size, board Board, excludedList []Board) (satisfiable bool, result Board) {
	formula := Reduce(n, board, excludedList)
	satisfiable, assignment := sat.SolveOnce(formula)
	if satisfiable {
		result = NewBoard(n)
		for vi, boolean := range assignment {
			if boolean {
				pi := v2p[n][vi]
				result[pi.row][pi.col] = pi.val
			}
		}
	}
	return satisfiable, result
}

// Result :
type Result struct {
	n            int
	board        Board
	excludedList []Board
}

// Next :
func (result *Result) Next() (satisfiable bool, board Board) {
	satisfiable, board = SolveOnce(result.n, result.board, result.excludedList)
	result.excludedList = append(result.excludedList, board)
	return satisfiable, board
}

// SolveAll :
func SolveAll(n Size, board Board) (result *Result) {
	return &Result{
		n:            n,
		board:        board,
		excludedList: make([]Board, 0),
	}
}
