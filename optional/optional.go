package optional

// Optional type is used to represent a value that may or may not exist.
// This type compensates for absence of generic NIL in go.
type Optional[T any] struct {
	value T
	isSet bool
}

// Creates a new Optional[T] with a value.
func Just[T any](t T) *Optional[T] {
	return &Optional[T]{
		value: t,
		isSet: true,
	}
}

// Creates a new Optional[T] without a value.
func Nothing[T any]() *Optional[T] {
	return &Optional[T]{
		isSet: false,
	}
}

// Whether the optional value is set.
func (o *Optional[T]) IsSet() bool {
	return o.isSet
}

// The actual value of the optional. Use with care by
// checking existence of the value by IsSet() first.
func (o *Optional[T]) Value() T {
	return o.value
}

// Then applies f to the value of an Optional[T] if it is set.
// Otherwise does nothing.
func (o *Optional[T]) Then(f func(T) T) *Optional[T] {
	if o.isSet {
		o.value = f(o.value)
	}
	return o
}
