from sudoku.board.generator import generate
from sudoku.board.solver import board_to_formula, implication_once
from sudoku.gui.web_gui import Game


s = Game()

s.run()
