import time
from typing import Any
import multiprocessing as mp

class Session:
    def __init__(self, timeout: int=60):
        self.timeout = timeout
        self.pool: dict[int, tuple[float, Any]] = {}
        self.running = mp.Value("b", True, lock=True)
        self.process = mp.Process(target=self._loop)
        self.process.start()

    def __del__(self):
        with self.running.get_lock():
            self.running.value = False
        self.process.join()
    def _loop(self):
        while True:
            with self.running.get_lock():
                if self.running.value == False:
                    break
                for key, (last_access, data) in self.pool.items():
                    if time.time() - last_access > self.timeout:
                        print(f"del key {key}")
                        self.pool.pop(key)
                        print(f"number of active users: {len(self.pool)}")
            time.sleep(self.timeout)

    def set(self, key: int, data: Any):
        with self.running.get_lock():
            print(f"set key {key}")
            self.pool[key] = (time.time(), data)
            print(f"number of active users: {len(self.pool)}")

    def get(self, key: int) -> Any:
        with self.running.get_lock():
            print(f"get key {key}")
            data = self.pool[key][1]
            self.pool[key] = (time.time(), data)
            print(f"number of active users: {len(self.pool)}")
            return data
