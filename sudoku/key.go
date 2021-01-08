package sudoku

import (
	"strconv"
)

// FromString :
func FromString(n int, key string) (current Board, ok bool) {
	if len(key) < n*n*n*n {
		ok = false
		return current, ok
	}
	current = NewBoard(n)
	count := 0
	var err error
	for row := 0; row < n*n; row++ {
		for col := 0; col < n*n; col++ {
			current[row][col], err = strconv.Atoi(key[count : count+1])
			if err != nil {
				ok = false
				return current, ok
			}
			count++
		}
	}
	ok = true
	return current, ok
}

// ToString :
func ToString(current Board) string {
	out := ""
	for _, row := range current {
		for _, cell := range row {
			out += strconv.Itoa(cell)
		}
	}
	return out
}
