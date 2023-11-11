package gophunc

// Set is a unique collection of elements.
// It uses the uniqueness of keys in Go maps.
type Set[T comparable] map[T]struct{}

// NewSet creates a new Set from an array.
// It does not ensure the order of elements.
func NewSet[T comparable](array []T) Set[T] {
	s := make(Set[T])
	for _, v := range array {
		s[v] = struct{}{}
	}
	return s
}

// Add adds an element to a Set.
func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

// Delete deletes an element from a Set.
func (s Set[T]) Delete(v T) {
	delete(s, v)
}

// Has checks if an element is in a Set.
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Keys returns all keys of a Set.
func (s Set[T]) Keys() []T {
	keys := make([]T, 0)
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}
