package sudoku

import "math/rand"

type c struct {
	row int
	col int
}

func getNumSolution(n int, board Board, cap int) int {
	numSolution := 0
	result := SolveAll(n, board)
	for numSolution < cap {
		sat, _ := result.Next()
		if !sat {
			break
		}
		numSolution++
	}
	return numSolution
}

func getPosList(n int) []p {
	pList := make([]p, 0, n*n*n*n*n*n)
	for row := 0; row < n*n; row++ {
		for col := 0; col < n*n; col++ {
			for val := 1; val <= n*n; val++ {
				pList = append(pList, p{
					row: row,
					col: col,
					val: val,
				})
			}
		}
	}
	return pList
}
func getCellList(n int) []c {
	cList := make([]c, 0, n*n*n*n)
	for row := 0; row < n*n; row++ {
		for col := 0; col < n*n; col++ {
			cList = append(cList, c{
				row: row,
				col: col,
			})
		}
	}
	return cList
}

// Generate :
func Generate(n int, seed int) Board {
	s := rand.NewSource(int64(seed))
	r := rand.New(s)
	board := NewBoard(n)
	cList := getCellList(n)
	// generate unique solution board
	r.Shuffle(len(cList), func(i int, j int) {
		cList[i], cList[j] = cList[j], cList[i]
	})
	for _, ci := range cList {
		if board[ci.row][ci.col] != 0 {
			continue
		}
		board[ci.row][ci.col] = 1 + r.Intn(n*n)
		numSolution := getNumSolution(n, board, 2)
		if numSolution == 1 {
			break
		}
		if numSolution == 0 {
			board[ci.row][ci.col] = 0
		}
	}
	// simplify board
	r.Shuffle(len(cList), func(i int, j int) {
		cList[i], cList[j] = cList[j], cList[i]
	})
	for _, ci := range cList {
		if board[ci.row][ci.col] == 0 {
			continue
		}
		val := board[ci.row][ci.col]
		board[ci.row][ci.col] = 0
		numSolution := getNumSolution(n, board, 2)
		if numSolution == 0 {
			panic("something wrong")
		}
		if numSolution > 1 {
			board[ci.row][ci.col] = val
		}
	}
	return board
}
