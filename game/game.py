from enum import Enum
from typing import Optional, Iterator
import multiprocessing as mp
import numpy as np
import pygame

import board

class Game:
    cell_size: int = (300, 300)
    screen: pygame.Surface
    waiting_panel_surf: pygame.Surface = pygame.image.load("./game/assets/waiting_panel.png")
    waiting_surf: pygame.Surface = pygame.image.load("./game/assets/waiting.png")
    playing_panel_surf: pygame.Surface = pygame.image.load("./game/assets/playing_panel.png")
    block_surf: pygame.Surface = pygame.image.load("./game/assets/block.png")
    current_surf: pygame.Surface = pygame.image.load("./game/assets/current.png")
    initial_surf: pygame.Surface = pygame.image.load("./game/assets/initial.png")
    violation_surf: pygame.Surface = pygame.image.load("./game/assets/violation.png")
    value_surf_list: list[pygame.Surface] = [pygame.image.load(f"./game/assets/{num}.png") for num in range(10)]
    youwin_panel_surf: pygame.Surface = pygame.image.load("./game/assets/youwin_panel.png")
    def __init__(self, cell_size: int = 60):
        self.cell_size = cell_size
        self.screen = pygame.display.set_mode(size=(9*cell_size, 11*cell_size))
        self.waiting_panel_surf = pygame.transform.scale(Game.waiting_panel_surf, (9*cell_size, 2*cell_size))
        self.waiting_surf = pygame.transform.scale(Game.waiting_surf, (9*cell_size, 9*cell_size))
        self.playing_panel_surf = pygame.transform.scale(Game.playing_panel_surf, (9*cell_size, 2*cell_size))
        self.block_surf = pygame.transform.scale(Game.block_surf, (3*cell_size, 3*cell_size))
        self.current_surf = pygame.transform.scale(Game.current_surf, (cell_size, cell_size))
        self.initial_surf = pygame.transform.scale(Game.initial_surf, (cell_size, cell_size))
        self.violation_surf = pygame.transform.scale(Game.violation_surf, (cell_size, cell_size))
        self.value_surf_list = [pygame.transform.scale(value_surf, (cell_size, cell_size)) for value_surf in Game.value_surf_list]
        self.youwin_panel_surf = pygame.transform.scale(Game.youwin_panel_surf, (9*cell_size, 2*cell_size))

    def get_value_surf(self, num: int) -> pygame.Surface:
        return self.value_surf_list[num]

    def pos_to_cell(self, pos: tuple[int, int]) -> tuple[int, int]:
        x, y = pos
        col = x // self.cell_size
        row = y // self.cell_size
        return row, col

    def cell_to_pos_tl(self, cell: tuple[int, int]) -> tuple[int, int]:
        row, col = cell
        y = row * self.cell_size
        x = col * self.cell_size
        return x, y

    @staticmethod
    def create_board(seed: int) -> tuple[np.ndarray, np.ndarray]:
        current_board = 1 + board.generate(seed)
        solution_board = 1 + next(board.solve_all(current_board - 1))
        return current_board, solution_board

    def loop(self, seed: int):
        np.random.seed(seed)
        class State(Enum):
            WAITING = 0
            PLAYING = 1
            QUIT = 2
            ENDED = 3
        state: State = State.WAITING
        ctx = mp.get_context("spawn")
        mp_running = ctx.Value("b", True, lock=True)
        board_queue = ctx.Queue()
        def create_board():
            while True:
                with mp_running.get_lock():
                    if not mp_running.value:
                        break
                if board_queue.empty():
                    board_queue.put(Game.create_board(np.random.randint(0, 2**32)))
        board_process = mp.Process(target=create_board)
        board_process.start()
        initial_board = None
        current_board = None
        solution_board = None
        current_cell = None
        while state != State.QUIT:
            # draw
            self.screen.fill((255, 255, 255))
            if state == State.WAITING:
                self.screen.blit(*self._blit_waiting())
                self.screen.blit(*self._blit_panel(self.waiting_panel_surf))
            if state == State.PLAYING:
                self.screen.blits(self._blit_board_list(initial_board, current_board, current_cell), doreturn=False)
                self.screen.blit(*self._blit_panel(self.playing_panel_surf))
            if state == State.ENDED:
                self.screen.blits(self._blit_board_list(initial_board, current_board, current_cell), doreturn=False)
                self.screen.blit(*self._blit_panel(self.youwin_panel_surf))
            pygame.display.flip()
            # event
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    state = State.QUIT
                    break
                if state == state.WAITING:
                    break
                if state == state.PLAYING:
                    if event.type == pygame.MOUSEBUTTONDOWN:
                        current_cell = self.pos_to_cell(event.pos)
                    if event.type == pygame.KEYDOWN:
                        if event.unicode in [str(num) for num in range(1, 10)]:
                            current_value = int(event.unicode)
                            if initial_board[current_cell] == 0:
                                current_board[current_cell] = current_value
                        if event.unicode in ["x", "X"]:
                            if initial_board[current_cell] == 0:
                                current_board[current_cell] = 0
                        if event.unicode in ["r", "R"]:
                            current_board = np.array(initial_board, copy=True)
                if state == state.ENDED:
                    if event.type == pygame.KEYDOWN:
                        if event.unicode in ["n", "N"]:
                            state = State.WAITING
                        if event.unicode in ["r", "R"]:
                            current_board = np.array(initial_board, copy=True)


            # backend
            if state == State.WAITING:
                initial_board, solution_board = board_queue.get()
                current_board = np.array(initial_board, copy=True)
                print(solution_board)
                state = State.PLAYING
            if state == State.PLAYING:
                if (current_board == solution_board).all():
                    state = State.ENDED
        pygame.quit()
        with mp_running.get_lock():
            mp_running.value = False
        board_process.join()

    @staticmethod
    def get_violation_cell(cell: tuple[int, int]) -> Iterator[tuple[int, int]]:
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
    def get_violation(current_board: np.ndarray) -> np.ndarray:
        violation_board = np.zeros((9, 9), dtype=bool)
        for row in range(9):
            for col in range(9):
                cell = row, col
                if violation_board[cell]:
                    continue
                value = current_board[cell]
                if value != 0:
                    for vcell in Game.get_violation_cell(cell):
                        if current_board[vcell] == value:
                            violation_board[cell] = True
                            violation_board[vcell] = True
        return violation_board

    def _blit_panel(self, panel: pygame.Surface) -> tuple[pygame.Surface, pygame.Rect]:
        return panel, pygame.Rect((0, 9*self.cell_size), (9*self.cell_size, 2*self.cell_size))

    def _blit_youwin(self) -> tuple[pygame.Surface, pygame.Rect]:
        return self.youwin_surf, pygame.Rect((0, 0), (9*self.cell_size, 9*self.cell_size))

    def _blit_waiting(self) -> tuple[pygame.Surface, pygame.Rect]:
        return self.waiting_surf, pygame.Rect((0, 0), (9*self.cell_size, 9*self.cell_size))

    def _blit_board_list(self, initial_board: np.ndarray, current_board: np.ndarray, current_cell: Optional[tuple[int, int]]) -> list[tuple[pygame.Surface, pygame.Rect]]:
        sequence = []
        # initial
        for row in range(9):
            for col in range(9):
                x, y = self.cell_to_pos_tl((row, col))
                if 1 <= initial_board[row, col] and initial_board[row, col] <= 9:
                    sequence.append(
                        (self.initial_surf, pygame.Rect((x, y), (self.cell_size, self.cell_size))),
                    )
        # block
        for row in range(0, 9, 3):
            for col in range(0, 9, 3):
                x, y = self.cell_to_pos_tl((row, col))
                sequence.append(
                    (self.block_surf, pygame.Rect((x, y), (3*self.cell_size, 3*self.cell_size))),
                )
        # values
        for row in range(9):
            for col in range(9):
                x, y = self.cell_to_pos_tl((row, col))
                value = current_board[row, col]
                sequence.append(
                    (self.value_surf_list[value], pygame.Rect((x, y), (self.cell_size, self.cell_size))),
                )
        # violation
        violation_board = Game.get_violation(current_board)
        for row in range(9):
            for col in range(9):
                if violation_board[row, col]:
                    x, y = self.cell_to_pos_tl((row, col))
                    sequence.append(
                        (self.violation_surf, pygame.Rect((x, y), (self.cell_size, self.cell_size))),
                    )

        # current cell
        if current_cell is not None:
            x, y = self.cell_to_pos_tl(current_cell)
            sequence.append(
                (self.current_surf, pygame.Rect((x, y), (self.cell_size, self.cell_size))),
            )
        return sequence