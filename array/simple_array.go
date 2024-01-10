package array

import "github.com/eicc27/Gophunc/optional"

func (t *TypedArray[T, U]) SimpleFilter(f func(T) bool) *TypedArray[T, U] {
	return t.Filter(func(t T, _ int, _ []T) bool {
		return f(t)
	})
}

func (t *TypedArray[T, U]) SimpleFilterIndex(f func(T) bool) *TypedArray[int, any] {
	return t.FilterIndex(func(t T, _ int, _ []T) bool {
		return f(t)
	})
}

func (t *TypedArray[T, U]) SimpleFlatMap(f func(T) []U) *TypedArray[U, any] {
	return t.FlatMap(func(t T, _ int, _ []T) []U {
		return f(t)
	})
}

func (t *TypedArray[T, U]) SimpleForEach(f func(T)) *TypedArray[T, U] {
	return t.ForEach(func(t T, _ int, _ []T) {
		f(t)
	})
}

// SimpleMap cuts off the filter function for reducing returning type of function from Optional[U] to U.
func (t *TypedArray[T, U]) SimpleMap(f func(T) U) *TypedArray[U, any] {
	return t.Map(func(t T, _ int, _ []T) *optional.Optional[U] {
		return optional.Just(f(t))
	})
}

// SimpleReduce asserts that the array has at least one of element without returning a potential error with Result.
func (t *TypedArray[T, U]) SimpleReduce(f func(T, T) T) T {
	r := t.Reduce(func(t1 T, t2 T, _ int, _ []T) T {
		return f(t1, t2)
	})
	return r.AsOK()
}
