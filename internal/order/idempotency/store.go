package idempotency

import "sync"

type Store[T any] struct {
	mu    sync.Mutex
	cache map[string]T
}

func New[T any]() *Store[T] {
	return &Store[T]{cache: make(map[string]T)}
}

func (s *Store[T]) Get(key string) (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.cache[key]
	return val, ok
}

func (s *Store[T]) Set(key string, val T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[key] = val
}
