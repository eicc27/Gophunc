package either

import O "github.com/eicc27/Gophunc/optional"

// Either is an option of two types.
// Either one type exists, or the other.
type Either[L any, R any] struct {
	Left  O.Optional[L]
	Right O.Optional[R]
}

// Creates a new Either[L, R] with a left value.
func Left[R any, L any](l L) *Either[L, R] {
	return &Either[L, R]{
		Left:  *O.Just(l),
		Right: *O.Nothing[R](),
	}
}

// Creates a new Either[L, R] with a right value.
func Right[L any, R any](r R) *Either[L, R] {
	return &Either[L, R]{
		Left:  *O.Nothing[L](),
		Right: *O.Just(r),
	}
}

// Flips right and left values for an Either.
func (e *Either[L, R]) Flip() *Either[R, L] {
	return &Either[R, L]{
		Left:  e.Right,
		Right: e.Left,
	}
}

// IsLeft checks if an Either[L, R] is a Left.
func (e *Either[L, R]) IsLeft() bool {
	return e.Left.IsSet()
}

// IsRight checks if an Either[L, R] is a Right.
func (e *Either[L, R]) IsRight() bool {
	return e.Right.IsSet()
}

// ToLeft regardless of the existence of left value, returns
// an O.Optional left value.
func (e *Either[L, R]) ToLeft() O.Optional[L] {
	return e.Left
}

// ToRight regardless of the existence of right value, returns
// an O.Optional right value.
func (e *Either[L, R]) ToRight() O.Optional[R] {
	return e.Right
}

// IfLeftThen applies f to the left value if the left value exists.
func (e *Either[L, R]) IfLeftThen(f func(L)) *Either[L, R] {
	if e.IsLeft() {
		f(e.Left.Value())
	}
	return e
}

// IfRightThen applies f to the right value of an Either[L, R] if it is a Right.
func (e *Either[L, R]) IfRightThen(f func(R)) *Either[L, R] {
	if e.IsRight() {
		f(e.Right.Value())
	}
	return e
}
