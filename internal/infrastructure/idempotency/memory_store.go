package idempotency

import (
	"sync"
	"time"
)

type Entry struct {
	Response interface{}
	Expires  time.Time
}

type Store struct {
	mu    sync.Mutex
	cache map[string]Entry
	ttl   time.Duration
}

func NewStore(ttl time.Duration) *Store {
	return &Store{
		cache: make(map[string]Entry),
		ttl:   ttl,
	}
}

func (s *Store) Get(key string) (interface{}, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.cache[key]
	if !ok || time.Now().After(entry.Expires) {
		delete(s.cache, key)
		return nil, false
	}

	return entry.Response, true
}

func (s *Store) Set(key string, response interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[key] = Entry{
		Response: response,
		Expires:  time.Now().Add(s.ttl),
	}
}
