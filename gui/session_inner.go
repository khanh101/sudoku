package gui

import (
	"sync"
	"time"
)

type entry struct {
	lastAccess time.Time
	data       interface{}
}
type session struct {
	timeout time.Duration
	pool    map[string]*entry
	mtx     sync.RWMutex
}

func newSession(timeout time.Duration) *session {
	s := &session{
		timeout: timeout,
		pool:    make(map[string]*entry),
	}
	go s.cleanLoop()
	return s
}

func (s *session) cleanLoop() {
	for {
		time.Sleep(s.timeout)
		s.mtx.Lock()
		delKeyList := make([]string, 0)
		for key, entry := range s.pool {
			if time.Since(entry.lastAccess) > s.timeout {
				delKeyList = append(delKeyList, key)
			}
		}
		for _, key := range delKeyList {
			delete(s.pool, key)
		}
		s.mtx.Unlock()
	}
}

func (s *session) set(key string, data interface{}) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.pool[key] = &entry{
		lastAccess: time.Now(),
		data:       data,
	}
}

func (s *session) get(key string) interface{} {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	e, ok := s.pool[key]
	if !ok {
		return nil
	}
	return e.data
}

func (s *session) numActiveKey() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return len(s.pool)
}
