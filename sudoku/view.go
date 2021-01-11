package sudoku

// GameView :
type GameView struct {
	CurrentString string     `json:"current_string"`
	YouWin        bool       `json:"youwin"`
	CurrentBoard  Board      `json:"current_board"`
	InitialMask   [][]bool   `json:"initial_mask"`
	ViolationMask [][]bool   `json:"violation_mask"`
	Pointer       CellView   `json:"pointer"`
	Explanation   []CellView `json:"explanation"`
	Message       string     `json:"message"`
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
	Val int `json:"val"`
}

// ImplicationView :
type ImplicationView struct {
	Row int             `json:"row"`
	Col int             `json:"col"`
	Val int             `json:"val"`
	Exp []PlacementView `json:"explanation"`
}
