package array

// Range behaves like Python range.
// start is included, end is excluded if step > 0,
// and vice versa if step < 0.
// Typically step could not be 0 for it will result in a dead loop.
// Range will try to set step = 1 instead.
func Range(start int, end int, step int) []int {
	result := make([]int, 0)
	if step == 0 {
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
	return New(Range(start, end, step)...)
}

// Count behaves like Python range(0, end, 1).
func Count(end int) []int {
	return Range(0, end, 1)
}

// TypedCount wraps the result of count into a TypedArray.
func TypedCount(end int) *TypedArray[int, any] {
	return TypedRange(0, end, 1)
}

// CountStep behaves like Python range(0, end, step).
func CountStep(end int, step int) []int {
	return Range(0, end, step)
}

// TypedCountStep wraps the result of countStep into a TypedArray.
func TypedCountStep(end int, step int) *TypedArray[int, any] {
	return TypedRange(0, end, step)
}
