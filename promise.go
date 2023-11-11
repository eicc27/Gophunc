package gophunc

import (
	"errors"
	"sync"
)

type Promise[T any] struct {
	result chan T
	err    error
}

// NewPromise creates a new Promise[T] from a function.
func NewPromise[T any](f func() (T, error)) *Promise[T] {
	p := &Promise[T]{
		result: make(chan T),
	}
	go func() {
		res, err := f()
		if err != nil {
			p.err = err
		} else {
			p.result <- res
		}
		close(p.result)
	}()
	return p
}

// Then applies successFn to the result of a Promise[T] if it is successful.
func (p *Promise[T]) Then(successFn func(T) (T, error)) *Promise[T] {
	res, _ := p.Await()
	return NewPromise[T](func() (T, error) {
		return successFn(res)
	})
}

// Catch applies failFn to the error of a Promise[T] if it is failed.
func (p *Promise[T]) Catch(failFn func(error) error) *Promise[T] {
	res, err := p.Await()
	return NewPromise[T](func() (T, error) {
		return res, failFn(err)
	})
}

// Await waits for the result of a Promise[T]. Note that await closes
// the channel of the Promise[T] after it is called.
func (p *Promise[T]) Await() (T, error) {
	res := <-p.result
	return res, p.err
}

// AwaitAll awaits for the results of multiple Promise[T]s.
// It returns a slice of results and a slice of errors.
func AwaitAll[T any](promises ...*Promise[T]) *Promise[[]T] {
	return NewPromise(func() ([]T, error) {
		var wg sync.WaitGroup
		res := make([]T, 0)
		var e error
		for _, promise := range promises {
			wg.Add(1)
			go func(p *Promise[T]) {
				defer wg.Done()
				result, err := p.Await()
				if err != nil {
					e = err
					return
				}
				res = append(res, result)
			}(promise)
		}
		wg.Wait()
		if (e != nil) {
			return nil, e
		}
		return res, nil
	})
}

// AwaitAny waits for the first successful Promise[T].
// If all promises fail, it returns an error.
func AwaitAny[T any](promises ...*Promise[T]) *Promise[Optional[T]] {
	return NewPromise(func() (Optional[T], error) {
		var wg sync.WaitGroup
		resultChan := make(chan Optional[T], 1)
		errCount := 0

		for _, promise := range promises {
			wg.Add(1)
			go func(p *Promise[T]) {
				defer wg.Done()
				if res, err := p.Await(); err == nil {
					resultChan <- *Just(res)
				} else {
					errCount++
				}
			}(promise)
		}

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		if result, ok := <-resultChan; ok {
			return result, nil
		}
		return *Nothing[T](), errors.New("all promises failed")
	})
}

// Await waits for the result of a Promise[T]. Note that await closes
// the channel of the Promise[T] after it is called.
func Await[T any](p *Promise[T]) (T, error) {
	res := <-p.result
	return res, p.err
}
