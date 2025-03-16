package tools

import "sync"

// BMap is the blocking Map
type BMap[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]V
}

func NewBMap[K comparable, V any]() *BMap[K, V] {
	return &BMap[K, V]{
		m: make(map[K]V),
	}
}

func (s *BMap[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *BMap[K, V]) Get(key K) (V, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.m[key]
	return val, ok
}
