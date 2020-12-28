package sudoku

import (
	"fmt"
	"strconv"
)

// Game :
type Game interface {
	View() GameView
	Place(PlacementView)
	Undo() (bool, PlacementView)
	Implication() (bool, PlacementView)
}

// GameView :
type GameView struct {
	YouWin        bool     `json:"youwin"`
	CurrentBoard  Board    `json:"current_board"`
	InitialMask   [][]bool `json:"initial_mask"`
	ViolationMask [][]bool `json:"violation_mask"`
}

// PlacementView :
type PlacementView struct {
	Row int `json:"row"`
	Col int `json:"col"`
	Val int `json:"value"`
}

func printBoard(board Board) {
	out := ""
	for _, row := range board {
		for _, val := range row {
			out += strconv.Itoa(val) + " "
		}
		out += "\n"
	}
	fmt.Println(out)
}

// NewGame :
func NewGame(n int, seed int) Game {
	current := Generate(n, seed)
	_, solution := SolveOnce(n, current, nil)
	initial := make([][]bool, n*n)
	for i := range initial {
		initial[i] = make([]bool, n*n)
	}
	for row := 0; row < n*n; row++ {
		for col := 0; col < n*n; col++ {
			initial[row][col] = current[row][col] != 0
		}
	}
	printBoard(solution)
	return &game{
		n:        n,
		current:  current,
		initial:  initial,
		solution: solution,
		stack:    make([]PlacementView, 0),
	}
}
