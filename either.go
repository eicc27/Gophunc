package gophunc

type Either[L any, R any] struct {
	left  Optional[L]
	right Optional[R]
}

// Creates a new Either[L, R] with a left value.
func Left[L any, R any](l L) *Either[L, R] {
	return &Either[L, R]{
		left:  *Just(l),
		right: *Nothing[R](),
	}
}

// Creates a new Either[L, R] with a right value.
func Right[L any, R any](r R) *Either[L, R] {
	return &Either[L, R]{
		left:  *Nothing[L](),
		right: *Just(r),
	}
}

// IsLeft checks if an Either[L, R] is a Left.
func (e *Either[L, R]) IsLeft() bool {
	return e.left.isSet
}

// IsRight checks if an Either[L, R] is a Right.
func (e *Either[L, R]) IsRight() bool {
	return e.right.isSet
}

// ThenIfLeft applies f to the left value of an Either[L, R] if it is a Left.
func (e *Either[L, R]) ThenIfLeft(f func(L) L) *Either[L, R] {
	if e.IsLeft() {
		e.left.value = f(e.left.value)
	}
	return e
}

// ThenIfRight applies f to the right value of an Either[L, R] if it is a Right.
func (e *Either[L, R]) ThenIfRight(f func(R) R) *Either[L, R] {
	if e.IsRight() {
		e.right.value = f(e.right.value)
	}
	return e
}
