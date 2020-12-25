from typing import Iterator

import numpy as np

from board.solver import solve_all

Pos = tuple[int, int, int]
Cell = tuple[int, int]

def get_num_solutions(board: np.ndarray, cap: int = 2) -> int:
    num_solutions = 0
    for _ in solve_all(board):
        num_solutions += 1
        if num_solutions >= cap:
            break
    return num_solutions

def get_position() -> Iterator[Pos]:
    position_list = []
    for y in range(9):
        for x in range(9):
            for v in range(9):
                pos = y, x, v
                position_list.append(pos)
    np.random.shuffle(position_list)
    return iter(position_list)

def get_cell() -> Iterator[Cell]:
    cell_list = []
    for y in range(9):
        for x in range(9):
            cell = y, x
            cell_list.append(cell)
    np.random.shuffle(cell_list)
    return iter(cell_list)

def generate(seed: int) -> np.ndarray:
    np.random.seed(seed)
    board = np.zeros((9, 9), dtype=int) - 1
    # generate unique solution board
    for pos in get_position():
        y, x, v = pos
        if 0 <= board[y, x] and board[y, x] < 9:
            continue
        board[y, x] = v
        num_solutions = get_num_solutions(board)
        if num_solutions == 1:
            break
        if num_solutions == 0:
            board[y, x] = -1
    # simplify board
    for cell in get_cell():
        y, x = cell
        if not (0 <= board[y, x] and board[y, x] < 9):
            continue
        v = board[y, x]
        board[y, x] = -1
        num_solutions = get_num_solutions(board)
        if num_solutions == 0:
            raise RuntimeError("something wrong")
        if num_solutions == 2:
            board[y, x] = v

    return board


