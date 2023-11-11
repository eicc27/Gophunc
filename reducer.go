package gophunc

// Reducer is a wrapper around a native array. It could carry out
// any generic operations on a native array, such as reduce, filter, foreach, etc.
// Most of its methods are designed chainable as long as the return does not
// necessarily break the type constraint.
// It supports interconversion with Mapper.
//
// See NewReducer for more details.
type Reducer[T any] struct {
	array []T
}

// NewReducer creates a new Reducer Array, which is a wrapper around a native array.
// Different from Mapper which has 2 type parameters, Reducer only has 1 type parameter, T,
// for reducer functions does not have to change the type of the elements.
func NewReducer[T any](array []T) *Reducer[T] {
	return &Reducer[T]{
		array: array,
	}
}

// Creates a new Reducer from a Mapper.
//
// See NewMapperFromReducer for the inverse operation.
func (m *Mapper[T, U]) NewReducerFromMapper() *Reducer[T] {
	return NewReducer(m.array)
}

// A typical ForEach implementation, chainable.
//
//	f: (item T, index int, array []T)
//
// f is applied for each element.
// Different from Map, ForEach does not return a new array.
// Instead, it applies f for each element and returns the reducer itself.
// This makes it more flexible as f can deal with outer variables.
func (r *Reducer[T]) ForEach(f func(T, int, []T)) *Reducer[T] {
	for i, v := range r.array {
		f(v, i, r.array)
	}
	return r
}

// A typical Reduce implementation.
//
//	f: (accumulator T, item T, index int, array []T) T
//
// Starting from the first element, accumulator is updated by applying f
// for each element. The final value of accumulator is returned.
func (r *Reducer[T]) Reduce(f func(T, T, int, []T) T) T {
	result := r.array[0]
	for i, v := range r.array[1:] {
		result = f(result, v, i, r.array)
	}
	return result
}

// Filter gets all elements that satisfy the predicate f. Chainable.
//
//	f: (item T, index int, array []T) bool
//
// For each element that is applied to f returns a true value, it is kept.
func (r *Reducer[T]) Filter(f func(T, int, []T) bool) *Reducer[T] {
	result := make([]T, 0)
	for i, v := range r.array {
		if f(v, i, r.array) {
			result = append(result, v)
		}
	}
	return NewReducer[T](result)
}

// FilterIndex gets all indices of elements that satisfy the predicate f.
// It is the indexed version of Filter, so it is not chainable.
func (r *Reducer[T]) FilterIndex(f func(T, int, []T) bool) []int {
	result := make([]int, 0)
	for i, v := range r.array {
		if f(v, i, r.array) {
			result = append(result, i)
		}
	}
	return result
}

// Splice does the operation in place, and returns the array of deleted elements.
// It takes the concept from JavaScript.
//
// Start could be negative, which means it is counted from the end.
//
// It first removes deleteCount elements starting from start(included),
// and then inserts items.
//
// If the deleteCount is more than the number of elements or negative,
// it removes all elements starting from start.
func (r *Reducer[T]) Splice(start int, deleteCount int, items ...T) *Reducer[T] {
	if start < 0 {
		start = len(r.array) + start
	}
	if start+deleteCount > len(r.array) || deleteCount < 0 {
		deleteCount = len(r.array) - start
	}
	deleted := r.array[start : start+deleteCount]
	r.array = append(r.array[:start], append(items, r.array[start+deleteCount:]...)...)
	return NewReducer(deleted)
}

// Slice takes the concept from JavaScript. It returns a new array.
// Start index is included, end index is excluded.
//
// Different from the basic Go implementation, it is chainable,
// and start and end both could take negative values.
//
// If start and end do not overlap, it returns an empty array.
func (r *Reducer[T]) Slice(start int, end int) *Reducer[T] {
	if start < 0 {
		start = len(r.array) + start
	}
	if end < 0 {
		end = len(r.array) + end
	}
	if start >= end {
		return NewReducer[T](make([]T, 0))
	}
	return NewReducer(r.array[start:end])
}
