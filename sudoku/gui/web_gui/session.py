import time
from typing import Any
import multiprocessing as mp

class Session:
    def __init__(self, timeout: int=60):
        self.pool = mp.Manager().dict()
        self.process = mp.Process(target=Session._loop, args=(self.pool, timeout))
        self.process.start()

    @staticmethod
    def _loop(pool, timeout):
        while True:
            for key in pool.keys():
                last_access, data = pool[key]
                if time.time() - last_access > timeout:
                    pool.pop(key)
                    print(f"del key: {key}, number of active users {len(pool)}")
            time.sleep(timeout)

    def set(self, key: int, data: Any):
        self.pool[key] = (time.time(), data)
        print(f"set key: {key}, number of active users {len(self.pool)}")

    def get(self, key: int) -> Any:
        print(f"get key {key}")
        data = self.pool[key][1]
        self.pool[key] = (time.time(), data)
        return data
