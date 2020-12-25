from typing import Iterator, Optional

import numpy as np

from sudoku.board.generator import generate
from sudoku.board.solver import solve_all
from sudoku.board.util import background


class Game:
    board_iterator: Iterator[tuple[np.ndarray, np.ndarray]]
    initial_board: Optional[np.ndarray]
    solution_board: Optional[np.ndarray]
    current_board: Optional[np.ndarray]

    def __init__(self, seed: int):
        self.board_iterator = background(Game.get_board_iterator(seed))
        self.initial_board = None
        self.solution_board = None
        self.current_board = None

    def __del__(self):
        del self.board_iterator

    def new_board(self):
        self.initial_board, self.solution_board = next(self.board_iterator)
        self.current_board = np.array(self.initial_board, copy=True)
        print(self.solution_board)

    def place(self, cell: tuple[int, int], value: int = 0):
        row, col = cell
        if 1 <= self.initial_board[row, col] and self.initial_board[row, col] <= 9:
            return
        self.current_board[row, col] = value

    def reset(self):
        self.current_board = np.array(self.initial_board, copy=True)

    def view(self) -> tuple[bool, np.ndarray, np.ndarray]:
        return (self.current_board == self.solution_board).all(), self.initial_board, self.current_board

    @staticmethod
    def get_board_iterator(seed: int) -> Iterator[tuple[np.ndarray, np.ndarray]]:
        np.random.seed(seed)
        while True:
            initial_board = 1 + generate(np.random.randint(0, 2 ** 32))
            solution_board = 1 + next(solve_all(initial_board - 1))
            yield initial_board, solution_board
