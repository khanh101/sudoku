package sudoku

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
