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
func NewPromise[T any](f func() Result[T]) *Promise[T] {
	p := &Promise[T]{
		result: make(chan T),
	}
	go func() {
		r := f()
		r.IfErrorThenApply(func(err error) {
			p.err = err
		}).IfOKThenApply(func(t T) {
			p.result <- t
		})
		close(p.result)
	}()
	return p
}

// Then applies successFn to the result of a Promise[T] if it is successful.
func (p *Promise[T]) Then(successFn func(T) Result[T]) *Promise[T] {
	return NewPromise[T](func() Result[T] {
		r := p.Await()
		var result Result[T]
		fn := func(t T) {
			result = successFn(t)
		}
		r.IfOKThenApply(fn)
		return result
	})
}

// Catch applies failFn to the error of a Promise[T] if it is failed.
func (p *Promise[T]) Catch(failFn func(error)) *Promise[T] {
	return NewPromise[T](func() Result[T] {
		r := p.Await()
		r.IfErrorThenApply(failFn)
		return r
	})
}

// Await waits for the result of a Promise[T]. Note that await closes
// the channel of the Promise[T] after it is called.
func (p *Promise[T]) Await() Result[T] {
	res := <-p.result
	return *NewResult(res, p.err)
}

// AwaitAll awaits for the results of multiple Promise[T]s.
// It returns a slice of results and a slice of errors.
func AwaitAll[T any](promises ...*Promise[T]) *Promise[[]T] {
	return NewPromise(func() Result[[]T] {
		var wg sync.WaitGroup
		res := make([]T, 0)
		var e error
		for _, promise := range promises {
			wg.Add(1)
			go func(p *Promise[T]) {
				defer wg.Done()
				r := p.Await()
				r.IfErrorThenApply(func(err error) {
					e = err
				}).IfOKThenApply(func(result T) {
					res = append(res, result)
				})
				if e != nil {
					return
				}
			}(promise)
		}
		wg.Wait()
		if e != nil {
			return *Error[[]T](e)
		}
		return *OK(res)
	})
}

// AwaitAny waits for the first successful Promise[T].
// If all promises fail, it returns an error.
func AwaitAny[T any](promises ...*Promise[T]) *Promise[Optional[T]] {
	return NewPromise(func() Result[Optional[T]] {
		var wg sync.WaitGroup
		resultChan := make(chan Optional[T], 1)
		errCount := 0

		for _, promise := range promises {
			wg.Add(1)
			go func(p *Promise[T]) {
				defer wg.Done()
				r := p.Await()
				r.IfOKThenApply(func(t T) {
					resultChan <- *Just(t)
				}).IfErrorThenApply(func(err error) {
					errCount++
				})
			}(promise)
		}

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		if result, ok := <-resultChan; ok {
			return *OK(result)
		}
		return *Error[Optional[T]](errors.New("all promises failed"))
	})
}

// Await waits for the result of a Promise[T]. Note that await closes
// the channel of the Promise[T] after it is called.
func Await[T any](p *Promise[T]) Result[T] {
	res := <-p.result
	return *NewResult(res, p.err)
}
