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

	board := sudoku.Generate(3, 1234)

	result := sudoku.SolveAll(3, board)
	for {
		sat, res := result.Next()
		if !sat {
			break
		}
		printBoard(res)
	}
	s := gui.NewServer()
	s.Run(":8080")
	return
}
