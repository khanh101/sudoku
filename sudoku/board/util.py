import multiprocessing as mp
from typing import Iterator, Any

class background(Iterator):
    iterator: Iterator[Any]
    def __init__(self, iterator: Iterator[Any]):
        self.parent_conn, self.child_conn = mp.Pipe()
        self.process = mp.Process(target=background._background_routine, args=(iterator, self.child_conn))
        self.process.start()
        self.parent_conn.send(True)

    def __del__(self):
        self.parent_conn.send(False)
        self.process.join()

    def __iter__(self):
        return self

    def __next__(self):
        value = self.parent_conn.recv()
        self.parent_conn.send(True)
        return value

    @staticmethod
    def _background_routine(iterator: Iterator[Any], conn: mp.Pipe):
        while conn.recv():
            try:
                conn.send(next(iterator))
            except StopIteration:
                conn.send(StopIteration)
                break

'''
def background(iterator: Iterator[Any]) -> Iterator[Any]:
    def _background_iterator(iterator: Iterator[Any], connection: mp.Pipe):
        while connection.recv():
            try:
                connection.send(next(iterator))
            except StopIteration:
                connection.send(StopIteration)
                break

    parent_connection, child_connection = mp.Pipe()
    process = mp.Process(target=_background_iterator, args=(iterator, child_connection))
    process.start()

    parent_connection.send(True) # first value
    while True:
        value = parent_connection.recv()
        if value is StopIteration:
            break
        parent_connection.send(True) # signal next value
        yield value

    process.join()
'''