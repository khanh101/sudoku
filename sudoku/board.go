package sudoku

// Board :
type Board [][]int

// NewBoard :
func NewBoard(n int) Board {
	out := make([][]int, n*n)
	for i := range out {
		out[i] = make([]int, n*n)
	}
	return out
}

// Copy :
func (board Board) Copy() Board {
	out := make([][]int, len(board))
	for r, row := range board {
		out[r] = make([]int, len(row))
		copy(out[r], row)
	}
	return out
}
