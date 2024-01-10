package array

import (
	O "github.com/eicc27/Gophunc/optional"
	"github.com/eicc27/Gophunc/set"
)

type TypedMap[T comparable, U any] struct {
	m map[T]U
}

// NewTypedMap creates a new TypedMap.
func NewTypedMap[T comparable, U any]() *TypedMap[T, U] {
	return &TypedMap[T, U]{
		m: make(map[T]U),
	}
}

// NewTypedMapFrom creates a new TypedMap from an existing map.
func NewTypedMapFrom[T comparable, U any](m map[T]U) *TypedMap[T, U] {
	return &TypedMap[T, U]{
		m: m,
	}
}

// Get returns the value of the key.
func (m *TypedMap[T, U]) Get(key T) *O.Optional[U] {
	if v, ok := m.m[key]; ok {
		return O.Just(v)
	}
	return O.Nothing[U]()
}

// Set sets the value of the key.
// If the key does not exist, it will be created.
func (m *TypedMap[T, U]) Set(key T, value U) *TypedMap[T, U] {
	m.m[key] = value
	return m
}

// Delete deletes the key.
func (m *TypedMap[T, U]) Delete(key T) *TypedMap[T, U] {
	delete(m.m, key)
	return m
}

// Keys returns the keys of the map.
func (m *TypedMap[T, U]) Keys() *TypedArray[T, any] {
	keys := make([]T, 0)
	for k := range m.m {
		keys = append(keys, k)
	}
	return New(keys...)
}

// Values returns the values of the map.
func (m *TypedMap[T, U]) Values() *TypedArray[U, any] {
	values := make([]U, 0)
	for _, v := range m.m {
		values = append(values, v)
	}
	return New(values...)
}

// ForEach applies f to each key-value pair.
func (m *TypedMap[T, U]) ForEach(f func(T, U)) *TypedMap[T, U] {
	for k, v := range m.m {
		f(k, v)
	}
	return m
}

// ToSet converts the keys of the map to a set.
func (m *TypedMap[T, U]) ToSet() set.Set[T] {
	s := set.NewSet[T]()
	for k := range m.m {
		s.Add(k)
	}
	return s
}

// GroupBy groups the array by the key returned by f.
func GroupBy[K comparable, U, V any](a *TypedArray[U, V], f func(U, int, []U) K) *TypedMap[K, *TypedArray[U, V]] {
	m := NewTypedMap[K, *TypedArray[U, V]]()
	for i, v := range a.array {
		key := f(v, i, a.array)
		if !m.Get(key).IsSet() {
			m.Set(key, NewMapper[V](v))
		} else {
			m.Get(key).Value().Push(v)
		}
	}
	return m
}
