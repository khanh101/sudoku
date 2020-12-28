package main

import (
	"fmt"
	"strconv"

	"github.com/khanhhhh/sudoku/gui"
	"github.com/khanhhhh/sudoku/sudoku"
)

func printBoard(board sudoku.Board) {
	out := ""
	for _, row := range board {
		for _, val := range row {
			out += strconv.Itoa(val) + " "
		}
		out += "\n"
	}
	fmt.Println(out)
}

func main() {
	s := gui.NewServer()
	s.Run(":8080")
	return
}
