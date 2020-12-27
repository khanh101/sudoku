import time
from typing import Any

class Session:
    def __init__(self, timeout: float=60):
        self.timeout = timeout
        self.pool = {}

    def number_of_active_user(self) -> int:
        self._eliminate_inactive_user()
        length = len(self.pool)
        return length

    def _eliminate_inactive_user(self):
        key_set = set(self.pool.keys())
        for key in key_set:
            if key in self.pool:
                last_access, data = self.pool[key]
                if time.time() - last_access > self.timeout:
                    self.pool.pop(key)
                    print(f"del key: {key}, number of active users {len(self.pool)}")

    def set(self, key: int, data: Any):
        self._eliminate_inactive_user()
        self.pool[key] = (time.time(), data)
        print(f"set key: {key}, number of active users {len(self.pool)}")

    def get(self, key: int) -> Any:
        print(f"get key {key}")
        data = self.pool[key][1]
        self.pool[key] = (time.time(), data)
        return data
