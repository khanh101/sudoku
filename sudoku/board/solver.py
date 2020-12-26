from typing import Optional, Iterator

import numpy as np
from pysat.solvers import Glucose4

pos2var_map: np.ndarray = np.arange(1, 1 + 9 * 9 * 9).reshape((9, 9, 9))
var2pos_map: Optional[dict[int, tuple[int, int, int]]] = None

Var = int
Pos = tuple[int, int, int]  # row, col, value


def pos2var(pos: Pos) -> Var:
    return int(pos2var_map[pos])


def var2pos(var: Var) -> Pos:
    global var2pos_map
    if var2pos_map is None:
        var2pos_map = {}
        for row in range(9):
            for col in range(9):
                for value in range(9):
                    _pos = row, col, value
                    _var = pos2var(_pos)
                    var2pos_map[_var] = _pos
    return var2pos_map[var]


def solve_all(board: np.ndarray) -> Iterator[np.ndarray]:
    solution_list = []
    while True:
        solvable, solution = solve_once(board, solution_list)
        if not solvable:
            break
        solution_list.append(solution)
        yield solution


def board_to_formula(board: np.ndarray, excluded: Optional[list[np.ndarray]] = None) -> list[list[int]]:
    formula: list[list[int]] = []

    def add_unique_active_clause(pos_list: list[Pos]):
        var_list = [pos2var(pos) for pos in pos_list]
        # there is at least 1 active variable: a1 or a2 or ... or an
        formula.append([var for var in var_list])
        # there is no 2 active variables: (not a1 or not a2) and (not a1 or not a3)
        for var1idx in range(len(var_list)):
            var1 = var_list[var1idx]
            for var2idx in range(var1idx + 1, len(var_list)):
                var2 = var_list[var2idx]
                formula.append([-var1, -var2])

    def add_all_active_clause(pos_list: list[Pos]):
        var_list = [pos2var(pos) for pos in pos_list]
        for var in var_list:
            formula.append([var])

    def add_not_all_active_clause(pos_list: list[Pos]):
        var_list = [pos2var(pos) for pos in pos_list]
        formula.append([-var for var in var_list])  # not (a and b) = (not a or not b)

    # add board
    pos_list = []
    for row in range(9):
        for col in range(9):
            value = board[row, col]
            if 0 <= value < 9:
                pos = row, col, value
                pos_list.append(pos)
    add_all_active_clause(pos_list)
    # exclude board
    if excluded is not None:
        for excluded_board in excluded:
            pos_list = []
            for row in range(9):
                for col in range(9):
                    value = excluded_board[row, col]
                    pos = row, col, value
                    pos_list.append(pos)
            add_not_all_active_clause(pos_list)
    # each cell has only 1 value
    # each (row, col) has only 1 value
    for row in range(9):
        for col in range(9):
            pos_list = []
            for value in range(9):
                pos = row, col, value
                pos_list.append(pos)
            add_unique_active_clause(pos_list)
    # each (column, value) has only 1 row
    # each (col, value) has only 1 row
    for col in range(9):
        for value in range(9):
            pos_list = []
            for row in range(9):
                pos = row, col, value
                pos_list.append(pos)
            add_unique_active_clause(pos_list)
    # each (value, row) has only 1 column
    # each (value, row) has only 1 col
    for value in range(9):
        for row in range(9):
            pos_list = []
            for col in range(9):
                pos = row, col, value
                pos_list.append(pos)
            add_unique_active_clause(pos_list)
    # each (3x3 block, value) has 1 pos
    for tl_y in range(0, 9, 3):
        for tl_x in range(0, 9, 3):
            for value in range(9):
                pos_list = []
                for dy in range(3):
                    for dx in range(3):
                        pos = tl_y + dy, tl_x + dx, value
                        pos_list.append(pos)
                add_unique_active_clause(pos_list)

    return formula



def solve_once(board: np.ndarray, excluded: Optional[list[np.ndarray]] = None) -> tuple[bool, Optional[np.ndarray]]:
    formula = board_to_formula(board, excluded)
    # solve
    solver = Glucose4(bootstrap_with=formula)
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


def implication_once(board: np.ndarray) -> Optional[tuple[int, int, int]]:
    def sign(x):
        if x > 0:
            return +1
        if x < 0:
            return -1
        return 0

    formula = board_to_formula(board)
    truth = np.zeros(shape=(1 + 9 * 9 * 9), dtype=int) # -1: false, 0: unknown, +1: true
    # assign to all truth value
    for conj in formula:
        if len(conj) == 1:
            var = conj[0]
            if truth[abs(var)] == 0:
                truth[abs(var)] = sign(var)
    # detect unsat
    for conj in formula:
        var01 = [var for var in conj if sign(var) * truth[abs(var)] != -1]
        if len(var01) == 0:
            return None  # unsat
    # implication
    for conj in formula:
        var01 = [var for var in conj if sign(var) * truth[abs(var)] != -1]
        val01 = [sign(var) * truth[abs(var)] for var in var01]
        if 1 in val01:
            continue  # sat
        if len(var01) == 1:
            print(var01[0])
            if var01[0] > 0:
                return var2pos(var01[0])  # implication
    return None
