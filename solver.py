from typing import Optional, Iterator

import numpy as np
from pysat.solvers import Glucose4

pos2var_map: np.ndarray = np.arange(1, 1 + 9 * 9 * 9).reshape((9, 9, 9))
var2pos_map: Optional[dict[int, tuple[int, int, int]]] = None

Var = int
Pos = tuple[int, int, int]


def pos2var(pos: Pos) -> Var:
    return int(pos2var_map[pos])


def var2pos(var: Var) -> Pos:
    global var2pos_map
    if var2pos_map is None:
        var2pos_map = {}
        for y in range(9):
            for x in range(9):
                for v in range(9):
                    _pos = y, x, v
                    _var = pos2var(_pos)
                    var2pos_map[_var] = _pos
    return var2pos_map[var]


def solve_all(board: np.ndarray) -> Iterator[np.ndarray]:
    solution_list = []
    while True:
        solvable, solution = solve(board, solution_list)
        if not solvable:
            break
        solution_list.append(solution)
        yield solution


def solve(board: np.ndarray, excluded: Optional[list[np.ndarray]] = None) -> tuple[bool, Optional[np.ndarray]]:
    solver = Glucose4()

    def add_unique_active_clause(pos_list: list[Pos]):
        var_list = [pos2var(pos) for pos in pos_list]
        # there is at least 1 active variable: a1 or a2 or ... or an
        solver.add_clause([var for var in var_list])
        # there is no 2 active variables: (not a1 or not a2) and (not a1 or not a3)
        for var1idx in range(len(var_list)):
            var1 = var_list[var1idx]
            for var2idx in range(var1idx + 1, len(var_list)):
                var2 = var_list[var2idx]
                solver.add_clause([-var1, -var2])

    def add_all_active_clause(pos_list: list[Pos]):
        var_list = [pos2var(pos) for pos in pos_list]
        for var in var_list:
            solver.add_clause([var])

    def add_not_all_active_clause(pos_list: list[Pos]):
        var_list = [pos2var(pos) for pos in pos_list]
        solver.add_clause([-var for var in var_list])  # not (a and b) = (not a or not b)

    # add board
    pos_list = []
    for y in range(9):
        for x in range(9):
            v = board[y, x]
            if 0 <= v and v < 9:
                pos = y, x, v
                pos_list.append(pos)
    add_all_active_clause(pos_list)
    # exclude board
    if excluded is not None:
        for excluded_board in excluded:
            pos_list = []
            for y in range(9):
                for x in range(9):
                    v = excluded_board[y, x]
                    pos = y, x, v
                    pos_list.append(pos)
            add_not_all_active_clause(pos_list)
    # each cell has only 1 value
    # each (y, x) has only 1 v
    for y in range(9):
        for x in range(9):
            pos_list = []
            for v in range(9):
                pos = y, x, v
                pos_list.append(pos)
            add_unique_active_clause(pos_list)
    # each (column, value) has only 1 row
    # each (x, v) has only 1 y
    for x in range(9):
        for v in range(9):
            pos_list = []
            for y in range(9):
                pos = y, x, v
                pos_list.append(pos)
            add_unique_active_clause(pos_list)
    # each (value, row) has only 1 column
    # each (v, y) has only 1 x
    for v in range(9):
        for y in range(9):
            pos_list = []
            for x in range(9):
                pos = y, x, v
                pos_list.append(pos)
            add_unique_active_clause(pos_list)
    # each (3x3 block, value) has 1 pos
    for tl_y in range(0, 9, 3):
        for tl_x in range(0, 9, 3):
            for v in range(9):
                pos_list = []
                for dy in range(3):
                    for dx in range(3):
                        pos = tl_y + dy, tl_x + dx, v
                        pos_list.append(pos)
                add_unique_active_clause(pos_list)
    # solve
    solvable: bool = solver.solve()
    if not solvable:
        return False, None
    model: list[int] = solver.get_model()
    solution = np.empty(shape=(9, 9), dtype=int)
    for var in model:
        if var > 0:
            pos = var2pos(var)
            y, x, v = pos
            solution[y, x] = v
    return True, solution
