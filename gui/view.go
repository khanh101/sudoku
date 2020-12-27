package gui

type keyView struct {
	Key int `json:key`
}

type posView struct {
	Row int `json:row`
	Col int `json:col`
	Val int `json:value`
}
