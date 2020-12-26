from typing import Iterator, Optional, Any

import numpy as np

from sudoku.board.generator import generate
from sudoku.board.solver import solve_all

class View:
    youwin: bool
    current_board: np.ndarray
    initial_mask: np.ndarray
    violation_mask: np.ndarray
    def marshal(self) -> dict[str, Any]:
        return {
            "youwin": self.youwin,
            "current_board": self.current_board.tolist(),
            "initial_mask": self.initial_mask.tolist(),
            "violation_mask": self.violation_mask.tolist(),
        }

class Game:
    current_board: np.ndarray
    solution_board: np.ndarray
    initial_mask: np.ndarray

    def __init__(self, seed: int):
        self.current_board = 1 + generate(seed)
        self.solution_board = 1 + next(solve_all(self.current_board - 1))
        self.initial_mask = self.current_board.astype(bool)
        print(self.solution_board)

    def view(self) -> View:
        view = View()
        view.youwin = bool((self.solution_board == self.current_board).all())
        view.current_board = self.current_board
        view.initial_mask = self.initial_mask
        view.violation_mask = Game._get_violation_mask(self.current_board)
        return view

    def place(self, cell: tuple[int, int], value: int = 0):
        row, col = cell
        if self.initial_mask[row, col] == 0:
            self.current_board[row, col] = value

    def reset(self):
        self.current_board[~self.initial_mask] = 0

    @staticmethod
    def _get_violation_cell(cell: tuple[int, int]) -> Iterator[tuple[int, int]]:
        row, col = cell
        for r in range(9):
            if r != row:
                yield r, col
        for c in range(9):
            if c != col:
                yield row, c
        btlr = 3 * (row // 3)
        btlc = 3 * (col // 3)
        for r in range(3):
            for c in range(3):
                if btlr + r == row and btlc + c == col:
                    continue
                yield btlr + r, btlc + c

    @staticmethod
    def _get_violation_mask(current_board: np.ndarray) -> np.ndarray:
        violation_mask = np.zeros((9, 9), dtype=bool)
        for row in range(9):
            for col in range(9):
                cell = row, col
                if violation_mask[cell]:
                    continue
                value = current_board[cell]
                if value != 0:
                    for vcell in Game._get_violation_cell(cell):
                        if current_board[vcell] == value:
                            violation_mask[cell] = True
                            violation_mask[vcell] = True
        return violation_mask
