package sudoku

import (
	"sync"

	"github.com/khanhhhh/sudoku/sat"
)

type game struct {
	n         int
	current   Board
	initial   [][]bool
	solution  Board
	violation [][]bool
	pointer   CellView
	stack     []PlacementView
	mtx       sync.RWMutex
}

func (g *game) View() GameView {
	g.mtx.RLock()
	defer g.mtx.RUnlock()
	view := GameView{}
	view.YouWin = true
	for row := 0; row < g.n*g.n; row++ {
		for col := 0; col < g.n*g.n; col++ {
			if g.current[row][col] != g.solution[row][col] {
				view.YouWin = false
				break
			}
		}
	}
	view.CurrentBoard = g.current
	view.InitialMask = g.initial
	view.ViolationMask = g.violation
	view.Pointer = g.pointer
	return view
}

func (g *game) getViolationCell(cin c) []c {
	out := make([]c, 0, 3*g.n*g.n)
	for ri := 0; ri < g.n*g.n; ri++ {
		if ri != cin.row {
			out = append(out, c{
				row: ri,
				col: cin.col,
			})
		}
	}
	for ci := 0; ci < g.n*g.n; ci++ {
		if ci != cin.col {
			out = append(out, c{
				row: cin.row,
				col: ci,
			})
		}
	}
	btlr := g.n * (cin.row / g.n)
	btlc := g.n * (cin.col / g.n)
	for ri := 0; ri < g.n; ri++ {
		for ci := 0; ci < g.n; ci++ {
			row := btlr + ri
			col := btlc + ci
			if row != cin.row || col != cin.col {
				out = append(out, c{
					row: row,
					col: col,
				})

			}
		}
	}
	return out
}

func (g *game) getViolation() [][]bool {
	violation := make([][]bool, g.n*g.n)
	for i := range violation {
		violation[i] = make([]bool, g.n*g.n)
	}
	for row := 0; row < g.n*g.n; row++ {
		for col := 0; col < g.n*g.n; col++ {
			if violation[row][col] {
				continue
			}
			val := g.current[row][col]
			if val != 0 {
				for _, vc := range g.getViolationCell(c{row: row, col: col}) {
					if g.current[vc.row][vc.col] == val {
						violation[row][col] = true
						violation[vc.row][vc.col] = true
					}
				}
			}
		}
	}
	return violation
}

func (g *game) Implication() (ok bool, view ImplicationView) {
	g.mtx.RLock()
	defer g.mtx.RUnlock()
	formula := Reduce(g.n, g.current, nil)
	unsat, assignment, explanation := sat.Implication(formula, nil, true)
	// ok, assignment, explanation := Implication(formula, true)
	if unsat {
		ok = false
		return ok, view
	}
	pi2leaf := make(map[p][]sat.Literal)
	for vi, value := range assignment {
		if value == sat.ValueTrue {
			pi := v2p[g.n][vi]
			if g.current[pi.row][pi.col] == 0 {
				// generate explanation
				findLeaf := func(root sat.Literal) []sat.Literal {
					leaf := make([]sat.Literal, 0)
					stack := []sat.Literal{vi}
					visited := make(map[sat.Literal]struct{})
					var lcurrent sat.Literal

					for len(stack) > 0 {
						lcurrent, stack = stack[len(stack)-1], stack[:len(stack)-1]
						visited[lcurrent] = struct{}{}
						clauseIdx := explanation[lcurrent]
						if len(formula[clauseIdx]) == 1 { // leaf
							leaf = append(leaf, lcurrent)
						}
						for _, child := range formula[clauseIdx] {
							if child == lcurrent {
								continue
							}
							if _, ok := visited[-child]; !ok {
								stack = append(stack, -child)
							}
						}
					}
					return leaf
				}
				leaf := findLeaf(vi)
				pi2leaf[pi] = leaf
			}
		}
	}
	// find smallest explanation
	if len(pi2leaf) == 0 {
		ok = false
		return ok, view
	}
	var smallestLenLeaf int = 999999999
	var smallestLeaf []sat.Literal = nil
	var smallestPi p
	for pi, leaf := range pi2leaf {
		if len(leaf) < smallestLenLeaf {
			smallestLenLeaf = len(leaf)
			smallestLeaf = leaf
			smallestPi = pi
		}
	}

	view.Col = smallestPi.col
	view.Row = smallestPi.row
	view.Val = smallestPi.val
	view.Exp = make([]PlacementView, len(smallestLeaf))
	for i, l := range smallestLeaf {
		pos := v2p[g.n][abs(l)]
		view.Exp[i] = PlacementView{
			Row: pos.row,
			Col: pos.col,
			Val: pos.val,
		}
	}
	ok = true
	return ok, view
}

func (g *game) Undo() (ok bool, view PlacementView) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if len(g.stack) == 0 {
		return false, view
	}
	view = g.stack[len(g.stack)-1]
	g.stack = g.stack[:len(g.stack)-1]
	return true, view
}

func (g *game) Place(p PlacementView) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if g.validRowColVal(p.Row, p.Col, p.Val) {
		if !g.initial[p.Row][p.Col] {
			g.current[p.Row][p.Col] = p.Val
			if p.Val > 0 {
				g.stack = append(g.stack, p)
			}
			g.violation = g.getViolation()
		}
		g.pointer = CellView{
			Row: p.Row,
			Col: p.Col,
		}
	}
}

func (g *game) Point(pointer CellView) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if g.validRowCol(pointer.Row, pointer.Col) {
		g.pointer = pointer
	}
}

func (g *game) validRowCol(row int, col int) bool {
	return row >= 0 && row <= g.n*g.n && col >= 0 && col <= g.n*g.n
}
func (g *game) validRowColVal(row int, col int, val int) bool {
	return row >= 0 && row < g.n*g.n && col >= 0 && col < g.n*g.n && val >= 0 && val <= g.n*g.n
}

func abs(x int) int {
	if x > 0 {
		return +x
	}
	if x < 0 {
		return -x
	}
	return 0
}
