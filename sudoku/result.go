package sudoku

// Result :
type Result interface {
	Next() (satisfiable bool, board Board)
}

type result struct {
	n            int
	board        Board
	excludedList []Board
}

// Next :
func (result *result) Next() (satisfiable bool, board Board) {
	satisfiable, board = SolveOnce(result.n, result.board, result.excludedList)
	result.excludedList = append(result.excludedList, board)
	return satisfiable, board
}
