package main

import (
	"fmt"

	"github.com/khanhhhh/sudoku/sudokusolver"
)

func main() {
	board := sudokusolver.NewBoard(3)
	sat, result := sudokusolver.SolveOnce(3, board, nil)
	fmt.Println(board.ToString())
	fmt.Println(sat)
	fmt.Println(result.ToString())
}
