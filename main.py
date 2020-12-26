import time

from sudoku.gui import pygame_gui

g = pygame_gui.Game(int(time.time()) % 2 ** 32)
g.run()
del g
