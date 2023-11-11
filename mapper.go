package gophunc

// Mapper is a wrapper around a native array.
// It is used to carry out Map and FlatMap operations.
// For Map and FlatMap has a different type of output,
// due to the limitation of generics in Go,
// the output type must be specified when creating this array.
// Except explicit situations where you have to use Map and FlatMap,
// Reducer is more recommended for its simplicity.
//
// See NewMapper for more details.
type Mapper[T any, U any] struct {
	array []T
}

// NewMapper creates a new Mapper Array, which is a wrapper around a native array.
// It has 2 type parameters, T and U, which are the types of the input and output.
//
// Normally, you would only need to specify U for output value(s)
// for T is inferred from the input array.
func NewMapper[U any, T any](array []T) *Mapper[T, U] {
	return &Mapper[T, U]{
		array: array,
	}
}

// Creates a new Mapper from a Reducer, by specifying the output type U.
// For mechanism of generics in Go, it must be a top-level function,
// but not a method of Reducer.
//
// In contrast, as input type T can be inferred from mapper while creating
// new reducers, NewReducer is a method of Mapper.
func NewMapperFromReducer[T any, U any](r *Reducer[T]) *Mapper[T, U] {
	return &Mapper[T, U]{
		array: r.array,
	}
}

// Map function that combines Map and FilterMap.
//
//	f : (item T, index int, array []T) Optional[U]
//
// f is applied for each element and the result of every application
// is stored and returned, if the result is not empty.
// Returns a new array of type U.
func (m *Mapper[T, U]) Map(f func(T, int, []T) Optional[U]) *Reducer[U] {
	result := make([]U, 0)
	for i, v := range m.array {
		r := f(v, i, m.array)
		if !r.isSet {
			continue
		}
		result = append(result, r.value)
	}
	return NewReducer[U](result)
}

// FlatMap is a flatten version of Map.
//
//	f: (item T, index int, array []T) []U
//
// The result array of every application is flattened into a single array
// as the return value.
func (m *Mapper[T, U]) FlatMap(f func(T, int, []T) []U) *Reducer[U] {
	result := make([]U, 0)
	for i, v := range m.array {
		result = append(result, f(v, i, m.array)...)
	}
	return NewReducer[U](result)
}

// See ForEach in reducer.go.
func (m *Mapper[T, U]) ForEach(f func(T, int, []T)) *Mapper[T, U] {
	return NewMapperFromReducer[T, U](m.NewReducer().ForEach(f))
}

// See Reduce in reducer.go.
func (m *Mapper[T, U]) Reduce(f func(T, T, int, []T) T) T {
	return m.NewReducer().Reduce(f)
}

// See Filter in reducer.go.
func (m *Mapper[T, U]) Filter(f func(T, int, []T) bool) *Mapper[T, U] {
	return NewMapperFromReducer[T, U](m.NewReducer().Filter(f))
}

// See FilterIndex in reducer.go.
func (m *Mapper[T, U]) FilterIndex(f func(T, int, []T) bool) []int {
	return m.NewReducer().FilterIndex(f)
}

// See Splice in reducer.go.
func (m *Mapper[T, U]) Splice(start int, deleteCount int, items ...T) *Mapper[T, U] {
	return NewMapperFromReducer[T, U](m.NewReducer().Splice(start, deleteCount, items...))
}

// See Slice in reducer.go.
func (m *Mapper[T, U]) Slice(start int, end int) *Mapper[T, U] {
	return NewMapperFromReducer[T, U](m.NewReducer().Slice(start, end))
}
