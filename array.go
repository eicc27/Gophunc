package gophunc

// TypedArray is a wrapper around a native array.
// It is used to carry out array ops like Map, Reduce, ...
// For Map and FlatMap has a different type of output,
// due to the limitation of generics in Go,
// the output type somehow must be specified when creating this array.
//
// For the ease of usage, two types are derived:
//
// > TUArray, a typed array with both input and output types specified.
//
// > TArray, a typed array with only input type specified.
type TypedArray[T, U any] struct {
	array []T
}

// NewMapperArray creates a new Mapper Array.
// It has 2 type parameters, T and U, which are the types of the input and output.
//
// Normally, you would only need to specify U for output value(s)
// for T is inferred from the input array.
func NewMapperArray[U, T any](array []T) *TypedArray[T, U] {
	return &TypedArray[T, U]{
		array: array,
	}
}

// NewTypedArray creates a typed array without specifying output type U.
// It is used for actions except the Map family.
//
// If Map ops are required, use WithType to add a type to the array.
func NewTypedArray[T any](array []T) *TypedArray[T, any] {
	return &TypedArray[T, any]{
		array: array,
	}
}

// WithType adds an output type to a single-typed array.
func WithType[U, T any](t *TypedArray[T, any]) *TypedArray[T, U] {
	return &TypedArray[T, U]{
		array: t.array,
	}
}

// Map function that combines Map and FilterMap.
//
//	f : (item T, index int, array []T) Optional[U]
//
// f is applied for each element and the result of every application
// is stored and returned, if the result is not empty.
// Returns a new array of type U.
func (m *TypedArray[T, U]) Map(f func(T, int, []T) Optional[U]) *TypedArray[U, any] {
	result := make([]U, 0)
	for i, v := range m.array {
		r := f(v, i, m.array)
		if !r.isSet {
			continue
		}
		result = append(result, r.value)
	}
	return NewTypedArray(result)
}

// FlatMap is a flatten version of Map.
//
//	f: (item T, index int, array []T) []U
//
// The result array of every application is flattened into a single array
// as the return value.
func (m *TypedArray[T, U]) FlatMap(f func(T, int, []T) []U) *TypedArray[U, any] {
	result := make([]U, 0)
	for i, v := range m.array {
		result = append(result, f(v, i, m.array)...)
	}
	return NewTypedArray(result)
}

// A typical ForEach implementation, chainable.
//
//	f: (item T, index int, array []T)
//
// f is applied for each element.
// Different from Map, ForEach does not return a new array.
// Instead, it applies f for each element and returns the reducer itself.
// This makes it more flexible as f can deal with outer variables.
func (r *TypedArray[T, U]) ForEach(f func(T, int, []T)) *TypedArray[T, U] {
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
func (r *TypedArray[T, U]) Reduce(f func(T, T, int, []T) T) T {
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
func (r *TypedArray[T, U]) Filter(f func(T, int, []T) bool) *TypedArray[T, U] {
	result := make([]T, 0)
	for i, v := range r.array {
		if f(v, i, r.array) {
			result = append(result, v)
		}
	}
	return NewMapperArray[U](result)
}

// FilterIndex gets all indices of elements that satisfy the predicate f.
// It is the indexed version of Filter.
func (r *TypedArray[T, U]) FilterIndex(f func(T, int, []T) bool) *TypedArray[int, any] {
	result := make([]int, 0)
	for i, v := range r.array {
		if f(v, i, r.array) {
			result = append(result, i)
		}
	}
	return NewTypedArray(result)
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
func (r *TypedArray[T, U]) Splice(start int, deleteCount int, items ...T) *TypedArray[T, U] {
	if start < 0 {
		start = len(r.array) + start
	}
	if start+deleteCount > len(r.array) || deleteCount < 0 {
		deleteCount = len(r.array) - start
	}
	deleted := r.array[start : start+deleteCount]
	r.array = append(r.array[:start], append(items, r.array[start+deleteCount:]...)...)
	return NewMapperArray[U](deleted)
}

// Slice takes the concept from JavaScript. It returns a new array.
// Start index is included, end index is excluded.
//
// Different from the basic Go implementation, it is chainable,
// and start and end both could take negative values.
//
// If start and end do not overlap, it returns an empty array.
func (r *TypedArray[T, U]) Slice(start int, end int) *TypedArray[T, U] {
	if start < 0 {
		start = len(r.array) + start
	}
	if end < 0 {
		end = len(r.array) + end
	}
	if start >= end {
		return NewMapperArray[U](make([]T, 0))
	}
	return NewMapperArray[U](r.array[start:end])
}
