from enum import Enum
from typing import Optional, Iterator

import numpy as np
import pygame

from sudoku import board

pygame.init()


class Game:
    cell_size: int = 300
    waiting_panel_surf: pygame.Surface = pygame.image.load("./sudoku/pygame_gui/assets/waiting_panel.png")
    waiting_surf: pygame.Surface = pygame.image.load("./sudoku/pygame_gui/assets/waiting.png")
    playing_panel_surf: pygame.Surface = pygame.image.load("./sudoku/pygame_gui/assets/playing_panel.png")
    block_surf: pygame.Surface = pygame.image.load("./sudoku/pygame_gui/assets/block.png")
    current_surf: pygame.Surface = pygame.image.load("./sudoku/pygame_gui/assets/current.png")
    initial_surf: pygame.Surface = pygame.image.load("./sudoku/pygame_gui/assets/initial.png")
    violation_surf: pygame.Surface = pygame.image.load("./sudoku/pygame_gui/assets/violation.png")
    value_surf_list: list[pygame.Surface] = [pygame.image.load(f"./sudoku/pygame_gui/assets/{num}.png") for num in
                                             range(10)]
    youwin_panel_surf: pygame.Surface = pygame.image.load("./sudoku/pygame_gui/assets/youwin_panel.png")
    screen: pygame.Surface
    board_game: board.Game

    def __init__(self, seed: int, cell_size: int = 60):
        self.cell_size = cell_size
        self.waiting_panel_surf = pygame.transform.scale(Game.waiting_panel_surf, (9 * cell_size, 2 * cell_size))
        self.waiting_surf = pygame.transform.scale(Game.waiting_surf, (9 * cell_size, 9 * cell_size))
        self.playing_panel_surf = pygame.transform.scale(Game.playing_panel_surf, (9 * cell_size, 2 * cell_size))
        self.block_surf = pygame.transform.scale(Game.block_surf, (3 * cell_size, 3 * cell_size))
        self.current_surf = pygame.transform.scale(Game.current_surf, (cell_size, cell_size))
        self.initial_surf = pygame.transform.scale(Game.initial_surf, (cell_size, cell_size))
        self.violation_surf = pygame.transform.scale(Game.violation_surf, (cell_size, cell_size))
        self.value_surf_list = [pygame.transform.scale(value_surf, (cell_size, cell_size)) for value_surf in
                                Game.value_surf_list]
        self.youwin_panel_surf = pygame.transform.scale(Game.youwin_panel_surf, (9 * cell_size, 2 * cell_size))
        self.screen = pygame.display.set_mode(size=(9 * cell_size, 11 * cell_size))
        self.board_game = board.Game(seed)

    def __del__(self):
        del self.board_game

    def run(self):
        class State(Enum):
            WAITING = 0
            PLAYING = 1
            QUIT = 2
            ENDED = 3

        state = State.WAITING
        current_cell = None
        while state != State.QUIT:
            # view
            self.screen.fill((255, 255, 255))
            if state == State.WAITING:
                self.screen.blit(*self._blit_waiting())
                self.screen.blit(*self._blit_panel(self.waiting_panel_surf))
            if state == State.PLAYING:
                youwin, initial_board, current_board = self.board_game.view()
                self.screen.blits(self._blit_board_list(initial_board, current_board, current_cell), doreturn=False)
                self.screen.blit(*self._blit_panel(self.playing_panel_surf))
            if state == State.ENDED:
                youwin, initial_board, current_board = self.board_game.view()
                self.screen.blits(self._blit_board_list(initial_board, current_board, current_cell), doreturn=False)
                self.screen.blit(*self._blit_panel(self.youwin_panel_surf))
            pygame.display.flip()
            # controller
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    state = State.QUIT
                    break
                if state == state.WAITING:
                    break
                if state == state.PLAYING:
                    if event.type == pygame.MOUSEBUTTONDOWN:
                        current_cell = self._pos_to_cell(event.pos)
                    if event.type == pygame.KEYDOWN:
                        if event.unicode in [str(num) for num in range(1, 10)]:
                            value = int(event.unicode)
                            self.board_game.place(current_cell, value)
                        if event.unicode in ["x", "X"]:
                            self.board_game.place(current_cell, 0)
                        if event.unicode in ["r", "R"]:
                            self.board_game.reset()
                if state == state.ENDED:
                    if event.type == pygame.KEYDOWN:
                        if event.unicode in ["n", "N"]:
                            state = State.WAITING
                        if event.unicode in ["r", "R"]:
                            self.board_game.reset()

            if state == state.WAITING:
                self.board_game.new_board()
                state = State.PLAYING
            if state == State.PLAYING:
                youwin, initial_board, current_board = self.board_game.view()
                if youwin:
                    state = State.ENDED
        pygame.quit()

    def _pos_to_cell(self, pos: tuple[int, int]) -> Optional[tuple[int, int]]:
        x, y = pos
        col = x // self.cell_size
        row = y // self.cell_size
        if not((1 <= row < 9) and (1 <= col < 9)):
            return None
        return row, col

    def _cell_to_pos_tl(self, cell: tuple[int, int]) -> tuple[int, int]:
        row, col = cell
        y = row * self.cell_size
        x = col * self.cell_size
        return x, y

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
    def _get_violation_board(current_board: np.ndarray) -> np.ndarray:
        violation_board = np.zeros((9, 9), dtype=bool)
        for row in range(9):
            for col in range(9):
                cell = row, col
                if violation_board[cell]:
                    continue
                value = current_board[cell]
                if value != 0:
                    for vcell in Game._get_violation_cell(cell):
                        if current_board[vcell] == value:
                            violation_board[cell] = True
                            violation_board[vcell] = True
        return violation_board

    def _blit_panel(self, panel: pygame.Surface) -> tuple[pygame.Surface, pygame.Rect]:
        return panel, pygame.Rect((0, 9 * self.cell_size), (9 * self.cell_size, 2 * self.cell_size))

    def _blit_waiting(self) -> tuple[pygame.Surface, pygame.Rect]:
        return self.waiting_surf, pygame.Rect((0, 0), (9 * self.cell_size, 9 * self.cell_size))

    def _blit_board_list(self, initial_board: np.ndarray, current_board: np.ndarray,
                         current_cell: Optional[tuple[int, int]]) -> list[tuple[pygame.Surface, pygame.Rect]]:
        sequence = []
        # initial
        for row in range(9):
            for col in range(9):
                x, y = self._cell_to_pos_tl((row, col))
                if 1 <= initial_board[row, col] <= 9:
                    sequence.append(
                        (self.initial_surf, pygame.Rect((x, y), (self.cell_size, self.cell_size))),
                    )
        # block
        for row in range(0, 9, 3):
            for col in range(0, 9, 3):
                x, y = self._cell_to_pos_tl((row, col))
                sequence.append(
                    (self.block_surf, pygame.Rect((x, y), (3 * self.cell_size, 3 * self.cell_size))),
                )
        # values
        for row in range(9):
            for col in range(9):
                x, y = self._cell_to_pos_tl((row, col))
                value = current_board[row, col]
                sequence.append(
                    (self.value_surf_list[value], pygame.Rect((x, y), (self.cell_size, self.cell_size))),
                )
        # violation
        violation_board = Game._get_violation_board(current_board)
        for row in range(9):
            for col in range(9):
                if violation_board[row, col]:
                    x, y = self._cell_to_pos_tl((row, col))
                    sequence.append(
                        (self.violation_surf, pygame.Rect((x, y), (self.cell_size, self.cell_size))),
                    )

        # current cell
        if current_cell is not None:
            x, y = self._cell_to_pos_tl(current_cell)
            sequence.append(
                (self.current_surf, pygame.Rect((x, y), (self.cell_size, self.cell_size))),
            )
        return sequence
