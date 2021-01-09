package gui

// KeyView :
type KeyView struct {
	Key string `json:"key"`
}

// BoardView :
type BoardView struct {
	Board string `json:"board"`
}

// PosView :
type PosView struct {
	Key string `json:"key"`
	Row int    `json:"row"`
	Col int    `json:"col"`
	Val int    `json:"value"`
}
