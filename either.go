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
