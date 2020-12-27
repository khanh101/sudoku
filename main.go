package main

import (
	"fmt"
	"strconv"

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
	g := sudoku.NewGame(3, 1234)
	ok, view := g.Implication()
	fmt.Println(ok, view)
	return
}
