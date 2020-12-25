from game import Game
import time
g = Game()
g.loop(int(time.time()) % (2**32))