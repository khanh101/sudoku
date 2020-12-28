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
	seed := 1234
	s := gui.NewServer(seed)
	s.Run(":8080")
	return
}
