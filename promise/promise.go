package promise

import (
	"errors"
	"sync"

	R "github.com/eicc27/Gophunc/result"
)

// Ported from JavaScript and realized with channels and goroutines.
// Once a promise is constructed, the task starts as a goroutine immediately,
// and could not be interrupted or stopped from outside controls.
type Promise[T any] struct {
	fulfill chan T
	err     error
}

// New creates a new Promise[T] from a task function.
//
// Example:
//
// timer := func(i int) func() *result.Result[int] {
// 	return func() *result.Result[int] {
// 		time.Sleep(time.Duration(i) * time.Second)
// 		fmt.Println(i)
// 		return result.OK(i)
// 	}
// }
// task := func(j int, panic bool) *promise.Promise[int] {
// 	return promise.New(timer(j)).Then(func(i int) *result.Result[int] {
// 		timer(i + 1)()
// 		if panic {
// 			return result.Error[int](errors.New("Panic"));
// 		}
// 		return result.OK(i)
// 	})
// }
//	promise.All(task(1, false), task(2, false)).Await() // 1, 2(t_2), 2(t_1), _, 3
func New[T any](f func() *R.Result[T]) *Promise[T] {
	p := &Promise[T]{
		fulfill: make(chan T),
	}
	go func() {
		r := f()
		r.IfErrorThen(func(err error) {
			p.err = err
		}).IfOKThen(func(t T) {
			p.fulfill <- t
		})
		close(p.fulfill)
	}()
	return p
}

// Then applies successFn to the result of a Promise[T] if it is successful.
func (p *Promise[T]) Then(successFn func(T) *R.Result[T]) *Promise[T] {
	return New[T](func() *R.Result[T] {
		r := p.Await()
		var result R.Result[T]
		ok := func(t T) {
			result = *successFn(t)
		}
		fail := func(_ error) {
			result = *r
		}
		r.IfOKThen(ok).IfErrorThen(fail)
		return &result
	})
}

// Catch applies failFn to the error of a Promise[T] if it is failed.
func (p *Promise[T]) Catch(failFn func(error)) *Promise[T] {
	return New[T](func() *R.Result[T] {
		r := p.Await()
		r.IfErrorThen(failFn)
		return r
	})
}

// Await blocks the main goroutine and waits for the result of a Promise[T].
// Note that await closes the channel of the Promise[T] after it is called.
func (p *Promise[T]) Await() *R.Result[T] {
	res := <-p.fulfill
	return R.New(res, p.err)
}

// All awaits for the results of multiple Promise[T]s,
// no matter how the promise fulfills (success or error).
// It returns a slice of results, or a slice of errors.
func All[T any](promises ...*Promise[T]) *Promise[[]T] {
	return New(func() *R.Result[[]T] {
		// uses waitgroup to wait for all promises to complete
		var wg sync.WaitGroup
		res := make([]T, 0)
		errs := make([]error, 0)
		for _, promise := range promises {
			wg.Add(1)
			go func(p *Promise[T]) {
				defer wg.Done()
				r := p.Await()
				r.IfErrorThen(func(err error) {
					errs = append(errs, err)
				}).IfOKThen(func(result T) {
					res = append(res, result)
				})
			}(promise)
		}
		wg.Wait()
		if len(errs) != 0 {
			return R.Error[[]T](errors.Join(errs...))
		}
		return R.OK(res)
	})
}

// Any waits for the first successful Promise[T].
// If all promises fail, it returns an error.
func Any[T any](promises ...*Promise[T]) *Promise[T] {
	return New(func() *R.Result[T] {
		var wg sync.WaitGroup
		resultChan := make(chan T)

		for _, promise := range promises {
			wg.Add(1)
			go func(p *Promise[T]) {
				defer wg.Done()
				r := p.Await()
				r.IfOKThen(func(t T) {
					resultChan <- t
				})
			}(promise)
		}

		// If none of the promise returns successfully, this coroutine
		// would close eventually and the resultChan would not be
		// blocked. Any trial that reads the resultChan would result an
		// error. This is a clever way to detect whether all of the promises
		// fail.
		go func() {
			wg.Wait()
			close(resultChan)
		}()

		// if any of the promise returns successfully, it jumps
		// into the if statement. Otherwise, the goroutine above would
		// close the resultChan and this `ok` would be false.
		if result, ok := <-resultChan; ok {
			return R.OK(result)
		}
		return R.Error[T](errors.New("all promises failed"))
	})
}

// Await waits for the result of a Promise[T]. Note that await closes
// the channel of the Promise[T] after it is called.
func Await[T any](p *Promise[T]) *R.Result[T] {
	res := <-p.fulfill
	return R.New(res, p.err)
}
