package sudoku

import "github.com/khanh-nguyen-code/sudoku/sat"

var baseCNF = make(map[int]sat.CNF)
var v2p = make(map[int]map[v]p)
var p2v = make(map[int]map[p]v)

type v = int
type p struct {
	row int
	col int
	val int
}

// ReduceBase :
func ReduceBase(n int) sat.CNF {
	base, ok := baseCNF[n]
	if !ok {
		// baseCNF
		base = sat.NewCNF()
		// var2pos - pos2var
		v2p[n] = make(map[v]p)
		p2v[n] = make(map[p]v)
		var vi v = 0
		for row := 0; row < n*n; row++ {
			for col := 0; col < n*n; col++ {
				for val := 1; val <= n*n; val++ {
					vi++
					v2p[n][vi] = p{
						row: row,
						col: col,
						val: val,
					}
					p2v[n][p{
						row: row,
						col: col,
						val: val,
					}] = vi
				}
			}
		}
		// add contrainsts
		addUniqueActiveClause := func(pList []p) {
			vList := make([]v, len(pList))
			for i, pi := range pList {
				vList[i] = p2v[n][pi]
			}
			// at least 1 variable is active
			base = append(base, vList)
			// no pair active at the same time
			for i1 := range vList {
				for i2 := i1 + 1; i2 < len(vList); i2++ {
					base = append(base, []sat.Literal{-vList[i1], -vList[i2]})
				}
			}
		}
		// each (row, col) has only 1 val
		for row := 0; row < n*n; row++ {
			for col := 0; col < n*n; col++ {
				pList := make([]p, 0, n*n)
				for val := 1; val <= n*n; val++ {
					pList = append(pList, p{
						row: row,
						col: col,
						val: val,
					})
				}
				addUniqueActiveClause(pList)
			}
		}
		// each (col, val) has only 1 row
		for col := 0; col < n*n; col++ {
			for val := 1; val <= n*n; val++ {
				pList := make([]p, 0, n*n)
				for row := 0; row < n*n; row++ {
					pList = append(pList, p{
						row: row,
						col: col,
						val: val,
					})
				}
				addUniqueActiveClause(pList)
			}
		}
		// each (val, row) has only 1 col
		for val := 1; val <= n*n; val++ {
			for row := 0; row < n*n; row++ {
				pList := make([]p, 0, n*n)
				for col := 0; col < n*n; col++ {
					pList = append(pList, p{
						row: row,
						col: col,
						val: val,
					})
				}
				addUniqueActiveClause(pList)
			}
		}
		// each (nxn block, value) has 1 (row, col)
		for tlr := 0; tlr < n*n; tlr += n {
			for tlc := 0; tlc < n*n; tlc += n {
				for val := 1; val <= n*n; val++ {
					pList := make([]p, 0, n*n)
					for r := 0; r < n; r++ {
						for c := 0; c < n; c++ {
							pList = append(pList, p{
								row: tlr + r,
								col: tlc + c,
								val: val,
							})
						}
					}
					addUniqueActiveClause(pList)
				}
			}
		}
		// done
		baseCNF[n] = base
	}
	return base
}

// Reduce :
func Reduce(n int, board Board, excludedList []Board) sat.CNF {
	formula := ReduceBase(n).Copy()
	addAllActiveClause := func(pList []p) {
		vList := make([]v, len(pList))
		for i, pi := range pList {
			vList[i] = p2v[n][pi]
		}
		for _, vi := range vList {
			formula = append(formula, []sat.Literal{vi})
		}
	}
	addNotAllActiveClause := func(pList []p) {
		vList := make([]v, len(pList))
		for i, pi := range pList {
			vList[i] = p2v[n][pi]
		}
		lList := make([]sat.Literal, 0, len(vList))
		for _, vi := range vList {
			lList = append(lList, -vi)
		}
		formula = append(formula, lList)
	}
	// add board
	pList := make([]p, 0)
	for row := 0; row < n*n; row++ {
		for col := 0; col < n*n; col++ {
			if board[row][col] != 0 {
				pList = append(pList, p{
					row: row,
					col: col,
					val: board[row][col],
				})
			}
		}
	}
	addAllActiveClause(pList)
	// exclude board
	for _, exclude := range excludedList {
		pList := make([]p, 0)
		for row := 0; row < n*n; row++ {
			for col := 0; col < n*n; col++ {
				if exclude[row][col] != 0 {
					pList = append(pList, p{
						row: row,
						col: col,
						val: exclude[row][col],
					})
				}
			}
		}
		addNotAllActiveClause(pList)
	}
	return formula
}
