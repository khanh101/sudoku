package sudoku

import "github.com/khanhhhh/sudoku/sat"

type game struct {
	n        int
	current  Board
	initial  [][]bool
	solution Board
	stack    []PlacementView
}

func (g *game) View() GameView {
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
	view.ViolationMask = g.getViolation()
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

func (g *game) Implication() (ok bool, view PlacementView) {
	formula := Reduce(g.n, g.current, nil)
	unsat, assignment := sat.Implication(formula, nil)
	if unsat {
		ok = false
		return ok, view
	}
	for vi, value := range assignment {
		if value == sat.ValueTrue {
			pi := v2p[g.n][vi]
			if g.current[pi.row][pi.col] == 0 {
				ok = true
				view.Row = pi.row
				view.Col = pi.col
				view.Val = pi.val
				return ok, view
			}

		}
	}
	ok = false
	return ok, PlacementView{}
}

func (g *game) Undo() (ok bool, view PlacementView) {
	if len(g.stack) == 0 {
		return false, view
	}
	view = g.stack[len(g.stack)-1]
	g.stack = g.stack[:len(g.stack)-1]
	return true, view
}

func (g *game) Place(p PlacementView) {
	if !g.initial[p.Row][p.Col] {
		g.current[p.Row][p.Col] = p.Val
		if p.Val > 0 {
			g.stack = append(g.stack, p)
		}
	}
}
