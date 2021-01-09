package sudoku

// Game :
type Game interface {
	View() GameView
	Place(PlacementView)
	Undo() (bool, PlacementView)
	Implication() (bool, ImplicationView)
}

// GameView :
type GameView struct {
	YouWin        bool     `json:"youwin"`
	CurrentBoard  Board    `json:"current_board"`
	InitialMask   [][]bool `json:"initial_mask"`
	ViolationMask [][]bool `json:"violation_mask"`
}

// PlacementView :
type PlacementView struct {
	Row int `json:"row"`
	Col int `json:"col"`
	Val int `json:"value"`
}

// ImplicationView :
type ImplicationView struct {
	Row int             `json:"row"`
	Col int             `json:"col"`
	Val int             `json:"value"`
	Exp []PlacementView `json:"explanation"`
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
