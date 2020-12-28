package gui

// KeyView :
type KeyView struct {
	Key string `json:"key"`
}

// PosView :
type PosView struct {
	Key string `json:"key"`
	Row int    `json:"row"`
	Col int    `json:"col"`
	Val int    `json:"value"`
}
