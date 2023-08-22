package utils

import "golang.org/x/exp/maps"

type Hashable interface {
	HashCode() uint64
}

type Set[T Hashable] map[uint64]T

func (s *Set[T]) Values() []T {
	return maps.Values(*s)
}

func (s *Set[T]) Add(entry T) {
	(*s)[entry.HashCode()] = entry
}

func (s *Set[T]) Contains(entry T) bool {
	_, ok := (*s)[entry.HashCode()]
	return ok
}
