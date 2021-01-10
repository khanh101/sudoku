package sudoku

// Game :
type Game interface {
	View() GameView
	Place(value int)
	Undo() (bool, PlacementView)
	Implication() (bool, ImplicationView)
	Point(pointer CellView)
}

// NewGame :
func NewGame(n int, current Board) (Game, bool) {
	printBoard(current)
	ok, solution := SolveOnce(n, current, nil)
	if !ok {
		return nil, false
	}
	initial := make([][]bool, n*n)
	for i := range initial {
		initial[i] = make([]bool, n*n)
	}
	for row := 0; row < n*n; row++ {
		for col := 0; col < n*n; col++ {
			initial[row][col] = current[row][col] != 0
		}
	}
	printBoard(solution)
	g := &game{
		n:        n,
		current:  current,
		initial:  initial,
		solution: solution,
		stack:    make([]PlacementView, 0),
	}
	g.violation = g.getViolation()
	return g, true
}
