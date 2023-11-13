package array

import (
	"errors"

	O "github.com/eicc27/Gophunc/optional"
	R "github.com/eicc27/Gophunc/result"
)

// TypedArray is a wrapper around a native array.
// It is used to carry out array ops like Map, Reduce, ...
// For Map and FlatMap having a different type of output,
// due to the limitation of generics in Go,
// the output type somehow must be specified when creating this array.
type TypedArray[T, U any] struct {
	array []T
}

// NewMapperArray creates a new Mapper Array.
// It has 2 type parameters, T and U, which are the types of the input and output.
//
// Normally, you would only need to specify U for output value(s)
// for T is inferred from the input array.
func NewMapperArray[U, T any](items ...T) *TypedArray[T, U] {
	return &TypedArray[T, U]{
		array: items,
	}
}

// NewTypedArray craetes an array without specifying output type U.
// It is used for actions except the Map family.
//
// If Map ops are required, use WithType to add a type to the array.
func NewTypedArray[T any](items ...T) *TypedArray[T, any] {
	return NewMapperArray[any](items...)
}

// WithType adds an output type to a single-typed array.
// This leverages the single-typed array to input-output-typed array
// to execute Map and FlatMap.
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
//
// Example:
// 	a := array.NewMapperArray[int](1, 2, 3).Map(
// 		func(t1, i int, t2 []int) optional.Optional[int] {
// 			if (t1 == 1) {
// 				return *optional.Nothing[int]()
// 			}
// 			return *optional.Just(t1 + 1)
// 		},
// 	)
// 	fmt.Println(a) // 3, 4
func (m *TypedArray[T, U]) Map(f func(T, int, []T) O.Optional[U]) *TypedArray[U, any] {
	result := make([]U, 0)
	for i, v := range m.array {
		r := f(v, i, m.array)
		if !r.IsSet() {
			continue
		}
		result = append(result, r.Value())
	}
	return NewTypedArray(result...)
}

// FlatMap is a flatten version of Map.
//
//	f: (item T, index int, array []T) []U
//
// The result array of every application is flattened into a single array
// as the return value.
// 
// Example (Also see Range):
//  a := array.NewMapperArray[int](1, 2, 3).FlatMap(
//  	func(t1, i int, t2 []int) []int {
//  		return array.Range(0, t1, 1)
//  	},
//  )
//  fmt.Println(a)  // 0 0 1 0 1 2
func (m *TypedArray[T, U]) FlatMap(f func(T, int, []T) []U) *TypedArray[U, any] {
	result := make([]U, 0)
	for i, v := range m.array {
		result = append(result, f(v, i, m.array)...)
	}
	return NewTypedArray(result...)
}

// A typical ForEach implementation, chainable.
//
//	f: (item T, index int, array []T)
//
// f is applied for each element.
// Different from Map, ForEach does not return a new array.
// Instead, it applies f for each element and returns the reducer itself.
// This makes it more flexible as f can deal with outer variables.
// 
// Example:
//  b := array.NewTypedArray[int]()
//  	array.NewTypedArray(1, 2, 3).ForEach(
//  		func(t1, i int, t2 []int) {
//  			b.Push(t1 + 1)
//  		},
//  	)
//  fmt.Println(b) // 2 3 4
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
// 
// Example:
//  r := array.NewTypedArray(1, 2, 3).Reduce(
//  	func(t, t1, i int, t2 []int) int {
//  		return t + t1
//  	},
//  )
//  fmt.Println(r.Right.Value()) // 6
func (r *TypedArray[T, U]) Reduce(f func(T, T, int, []T) T) R.Result[T] {
	if r.Length() == 0 {
		return *R.Error[T](errors.New("array to reduce must have at leat 1 element"))
	}
	result := r.array[0]
	for i, v := range r.array[1:] {
		result = f(result, v, i, r.array)
	}
	return *R.OK(result)
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
	return NewMapperArray[U](result...)
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
	return NewTypedArray(result...)
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
// 
// If the start is too large(more than the length of the array),
// it will only do insertion at the end of the array(equals to push).
func (r *TypedArray[T, U]) Splice(start int, deleteCount int, items ...T) *TypedArray[T, U] {
	if (start >= len(r.array)) {
		r.Push(items...)
		return NewMapperArray[U, T]()
	}
	if start < 0 {
		start = len(r.array) + start
	}
	if start+deleteCount > len(r.array) || deleteCount < 0 {
		deleteCount = len(r.array) - start
	}
	deleted := r.Slice(start, start+deleteCount)
	r.array = append(r.array[:start], append(items, r.array[start+deleteCount:]...)...)
	return deleted
}

// Slice takes the concept from JavaScript. It returns a new array.
// Start index is included, end index is excluded.
//
// Different from the basic Go implementation, it is chainable,
// and start and end both could take negative values.
//
// If start and end do not overlap, or start is too large, it returns an empty array.
func (r *TypedArray[T, U]) Slice(start int, end int) *TypedArray[T, U] {
	if start >= len(r.array) {
		return NewMapperArray[U, T]()
	}
	if start < 0 {
		start = len(r.array) + start
	}
	if end < 0 {
		end = len(r.array) + end
	}
	if start >= end {
		return NewMapperArray[U, T]()
	}
	return NewMapperArray[U](r.array[start:end]...)
}

// Index the array with the given index.
// Supports negative index.
// If the index is too large, it will dropback to index = -1.
func (r *TypedArray[T, U]) At(index int) T {
	if index >= len(r.array) {
		index = -1
	}
	if index < 0 {
		index = len(r.array) + index
	}
	return r.array[index]
}

// Returns the length of the array.
func (r *TypedArray[T, U]) Length() int {
	return len(r.array)
}

// Push pushes some items at the end of the array.
func (r *TypedArray[T, U]) Push(items ...T) *TypedArray[T, U] {
	r.array = append(r.array, items...)
	return r
}

// Pop pops an item at the end of the array.
// If the array does not have any item to pop,
// it does nothing and returns a nothing optional.
func (r *TypedArray[T, U]) Pop() O.Optional[T] {
	if r.Length() < 1 {
		return *O.Nothing[T]()
	}
	popped := r.At(-1)
	*r = *r.Slice(0, -1)
	return *O.Just(popped)
}

// Shift pops an element at the first of the array.
// If the array does not have any item to pop,
// it does nothing and returns a nothing optional.
func (r *TypedArray[T, U]) Shift() O.Optional[T] {
	if r.Length() < 1 {
		return *O.Nothing[T]()
	}
	shifted := r.At(0)
	*r = *r.Slice(1, r.Length())
	return *O.Just(shifted)
}

// Unshift pushes items at the beginning of the array.
func (r *TypedArray[T, U]) Unshift(items ...T) *TypedArray[T, U] {
	r.array = append(items, r.array...)
	return r
}

// Range behaves like Python range.
// start is included, end is excluded if step > 0,
// and vice versa if step < 0.
// Typically step could not be 0 for it will result in a dead loop.
// Range will try to set step = 1 instead.
func Range(start int, end int, step int) []int {
	result := make([]int, 0)
	if (step == 0) {
		step = 1
	}
	if step > 0 {
		for i := start; i < end; i += step {
			result = append(result, i)
		}
	} else {
		for i := end; i > start; i -= step {
			result = append(result, i)
		}
	}
	return result
}

// TypedRange wraps the result of range into a TypedArray.
func TypedRange(start int, end int, step int) *TypedArray[int, any] {
	return NewTypedArray(Range(start, end, step)...)
}
