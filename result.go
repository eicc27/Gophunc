package gophunc

type Result[T any] Either[error, T]

// Creates a new OK result.
func OK[T any](t T) *Result[T] {
	return &Result[T]{
		left:  *Nothing[error](),
		right: *Just(t),
	}
}

// Craetes an error result.
func Error[T any](e error) *Result[T] {
	return &Result[T]{
		left:  *Just(e),
		right: *Nothing[T](),
	}
}

// Creates a new result with a result and a potential error.
//
// Example:
//
//	 info, err := os.Stat("go.mod")
//		result := NewResult(info, err)
//		result.IfOKThenApply(func(t fs.FileInfo) {
//			fmt.Println(t.Name(), t.Size())
//		}).IfErrorThenApply(func(err error) {
//			fmt.Println(err.Error())
//		})
func NewResult[T any](result T, e error) *Result[T] {
	if e != nil {
		return Error[T](e)
	}
	return OK(result)
}

func (r *Result[T]) IsOK() bool {
	return r.right.isSet
}

func (r *Result[T]) IsError() bool {
	return !r.IsOK()
}

func (r *Result[T]) IfOKThen(f func(T) T) *Result[T] {
	if r.IsOK() {
		r.right = *Just(f(r.right.value))
	}
	return r
}

func (r *Result[T]) IfErrorThen(f func(error) error) *Result[T] {
	if r.IsError() {
		r.left = *Just(f(r.left.value))
	}
	return r
}

func (r *Result[T]) IfOKThenApply(f func(T)) *Result[T] {
	if r.IsOK() {
		f(r.right.value)
	}
	return r
}

func (r *Result[T]) IfErrorThenApply(f func(error)) *Result[T] {
	if r.IsError() {
		f(r.left.value)
	}
	return r
}
