package gophunc

// Either is an option of two types.
// Either one type exists, or the other.
// Chaining an either type is possible to create a multiple choice of types.
type Either[L any, R any] struct {
	left  Optional[L]
	right Optional[R]
}

// Creates a new Either[L, R] with a left value.
func Left[R any, L any](l L) *Either[L, R] {
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

// Flips right and values for an Either.
func (e *Either[L, R]) Flip() *Either[R, L] {
	return &Either[R, L]{
		left:  e.right,
		right: e.left,
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

// ToLeft regardless of the existence of left value, returns
// an Optional left value.
func (e *Either[L, R]) ToLeft() Optional[L] {
	return e.left
}

// ToRight regardless of the existence of right value, returns
// an Optional right value.
func (e *Either[L, R]) ToRight() Optional[R] {
	return e.right
}

// IfLeftThen applies f to the left value of an Either[L, R] if it is a Left.
func (e *Either[L, R]) IfLeftThen(f func(L) L) *Either[L, R] {
	if e.IsLeft() {
		e.left.value = f(e.left.value)
	}
	return e
}

func (e *Either[L, R]) IfLeftThenApply(f func(L)) *Either[L, R] {
	if e.IsLeft() {
		f(e.left.value)
	}
	return e
}

// IfRightThen applies f to the right value of an Either[L, R] if it is a Right.
func (e *Either[L, R]) IfRightThen(f func(R) R) *Either[L, R] {
	if e.IsRight() {
		e.right.value = f(e.right.value)
	}
	return e
}

func (e *Either[L, R]) IfRightThenApply(f func(R)) *Either[L, R] {
	if e.IsRight() {
		f(e.right.value)
	}
	return e
}
