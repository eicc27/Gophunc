package result

import (
	E "github.com/eicc27/Gophunc/either"
	O "github.com/eicc27/Gophunc/optional"
)

// A type derived from Either[L, R]. The type L is 
// locked as an error type to facilitate error handling.
type Result[T any] E.Either[error, T]

// Creates a new OK result.
func OK[T any](t T) *Result[T] {
	return &Result[T]{
		Left: *O.Nothing[error](),
		Right: *O.Just(t),
	}
}

// Craetes an error result.
func Error[T any](e error) *Result[T] {
	return &Result[T]{
		Left:  *O.Just(e),
		Right: *O.Nothing[T](),
	}
}

// Creates a new result with a result and a potential error.
//
// Example:
//
//	info, err := os.Stat("go.mod")
//   result := result.NewResult(info, err)
//   result.IfOKThen(func(t fs.FileInfo) {
//   	fmt.Println(t.Name(), t.Size())
//   }).IfErrorThen(func(err error) {
//   	fmt.Println(err.Error())
//   })
func NewResult[T any](result T, e error) *Result[T] {
	if e != nil {
		return Error[T](e)
	}
	return OK(result)
}

// Checks whether this result is OK.
func (r *Result[T]) IsOK() bool {
	return r.Right.IsSet()
}

// Checks whether this result is an error.
func (r *Result[T]) IsError() bool {
	return !r.IsOK()
}

// Applies a function to the value if the result is OK.
// Otherwise does nothing.
func (r *Result[T]) IfOKThen(f func(T)) *Result[T] {
	if r.IsOK() {
		f(r.Right.Value())
	}
	return r
}

// Applies a function to the error. If the result is OK
// does nothing.
func (r *Result[T]) IfErrorThen(f func(error)) *Result[T] {
	if r.IsError() {
		f(r.Left.Value())
	}
	return r
}
