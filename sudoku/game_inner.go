package sudoku

import (
	"fmt"
	"sync"

	"github.com/khanhhhh/sudoku/sat"
)

type game struct {
	n           int
	current     Board
	initial     [][]bool
	solution    Board
	violation   [][]bool
	pointer     CellView
	stack       []CellView
	explanation []CellView
	message     string
	mtx         sync.RWMutex
}

func (g *game) View() GameView {
	return g.rlock(func() interface{} {
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
		view.CurrentString = ToString(g.current)
		view.CurrentBoard = g.current
		view.InitialMask = g.initial
		view.ViolationMask = g.violation
		view.Pointer = g.pointer
		view.Message = g.message
		if g.explanation != nil {
			view.Explanation = g.explanation
		} else {
			view.Explanation = make([]CellView, 0)
		}
		return view
	}).(GameView)
}

func (g *game) Implication() bool {
	formula := g.rlock(func() interface{} {
		return Reduce(g.n, g.current, nil)
	}).(sat.CNF)
	unsat, assignment, explanation := sat.Implication(formula, nil, true)
	if unsat {
		g.lock(func() interface{} {
			g.message = "board is unsatisfiable"
			return nil
		})
		return false
	}
	pi2leaf := g.rlock(func() interface{} {
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
		return pi2leaf
	}).(map[p][]sat.Literal)
	// find smallest explanation
	if len(pi2leaf) == 0 {
		g.lock(func() interface{} {
			g.message = "implication not found"
			return nil
		})
		return false
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
	///
	g.Point(CellView{
		Row: smallestPi.row,
		Col: smallestPi.col,
	})
	g.Place(smallestPi.val)
	g.lock(func() interface{} {
		g.message = fmt.Sprintf("implication found at {row: %d, col: %d, val %d}", smallestPi.row, smallestPi.col, smallestPi.val)
		g.explanation = make([]CellView, len(smallestLeaf))
		for i, l := range smallestLeaf {
			pos := v2p[g.n][abs(l)]
			g.explanation[i] = CellView{
				Row: pos.row,
				Col: pos.col,
			}
		}
		return nil
	})
	return true
}

func (g *game) Undo() bool {
	if g.rlock(func() interface{} {
		return len(g.stack) == 0
	}).(bool) {
		return false
	}
	for {
		value := g.lock(func() interface{} {
			if len(g.stack) == 0 {
				return nil
			}
			action := g.stack[len(g.stack)-1]
			g.stack = g.stack[:len(g.stack)-1]
			return action
		})
		if value == nil {
			g.lock(func() interface{} {
				g.message = "stack empty"
				return nil
			})
			return false
		}
		action := value.(CellView)
		g.Point(CellView{
			Row: action.Row,
			Col: action.Col,
		})
		if g.Place(0) {
			g.lock(func() interface{} {
				g.message = fmt.Sprintf("undo found {row: %d, col: %d}", action.Row, action.Col)
				return nil
			})
			return true
		}
	}
}

func (g *game) Place(v int) bool {
	return g.lock(func() interface{} {
		r, c := g.pointer.Row, g.pointer.Col
		if g.initial[r][c] {
			return false
		}
		if g.current[r][c] == v {
			return false
		}
		g.current[r][c] = v
		g.explanation = nil
		if v != 0 {
			g.stack = append(g.stack, CellView{
				Row: r,
				Col: c,
			})
		}
		g.violation = g.getViolation()
		return true
	}).(bool)
}

func (g *game) Point(pointer CellView) {
	g.lock(func() interface{} {
		if g.validRowCol(pointer.Row, pointer.Col) {
			g.pointer = pointer
		}
		return nil
	})
}
