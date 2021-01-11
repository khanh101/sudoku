package gui

// KeyView :
type KeyView struct {
	Key string `json:"key"`
}

// BoardView :
type BoardView struct {
	Board string `json:"board"`
}

// PlaceView :
type PlaceView struct {
	Key string `json:"key"`
	Val int    `json:"val"`
}

// PointView :
type PointView struct {
	Key string `json:"key"`
	Row int    `json:"row"`
	Col int    `json:"col"`
}
