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
	board := sudoku.NewBoard(3)
	board[0][0] = 1
	result := sudoku.SolveAll(3, board)
	for {
		sat, board := result.Next()
		if sat {
			printBoard(board)
		}
	}
}
