package main

import (
	"github.com/khanhhhh/sudoku/gui"
)

func main() {
	s := gui.NewServer()
	s.Run(":8080")
	return
}
