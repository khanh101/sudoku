from sudoku.board.generator import generate
from sudoku.board.solver import board_to_formula, implication_once
from sudoku.gui.web_gui import Game

a = generate(1231)
b = board_to_formula(a)

print(implication_once(a))



s = Game()

s.run()
