package sudoku

// GameView :
type GameView struct {
	YouWin        bool     `json:"youwin"`
	CurrentBoard  Board    `json:"current_board"`
	InitialMask   [][]bool `json:"initial_mask"`
	ViolationMask [][]bool `json:"violation_mask"`
	Pointer       CellView `json:"pointer"`
}

// CellView :
type CellView struct {
	Row int `json:"row"`
	Col int `json:"col"`
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
