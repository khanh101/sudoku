import time
from typing import Any
import multiprocessing as mp

class Session:
    def __init__(self, timeout: int=60):
        self.timeout = timeout
        self.pool: dict[int, tuple[float, Any]] = {}
        self.running = mp.Value("b", True, lock=True)
        self.process = mp.Process(target=self._loop)

    def __del__(self):
        with self.running.get_lock():
            self.running.value = False
    def _loop(self):
        while True:
            with self.running.get_lock():
                if self.running.value == False:
                    break
                for key, (last_access, data) in self.pool.items():
                    if time.time() - last_access > self.timeout:
                        print(f"del key {key}")
                        self.pool.pop(key)
            time.sleep(self.timeout)

    def set(self, key: int, data: Any):
        print(f"set key {key}")
        with self.running.get_lock():
            self.pool[key] = (time.time(), data)

    def get(self, key: int) -> Any:
        print(f"get key {key}")
        with self.running.get_lock():
            data = self.pool[key][1]
            self.pool[key] = (time.time(), data)
            return data
