import numpy as np

from sudoku.board.generator import generate
from sudoku.board.solver import board_to_formula, implication_once
from sudoku.gui.web_gui import Game

a = generate(1231)

board = -np.ones((9, 9), dtype=int)
for i in range(8):
    board[0, i] = i
print(board)
print(implication_once(board))



s = Game()

s.run()
