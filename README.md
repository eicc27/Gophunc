# Gophunc

Functional collections and monads for go. Mapper & Reducer, Promise, Set, Either & Optional.

It needs **Go 1.18+** to provide generic types.


## Mapper

`Mapper[T, U]` takes two generic types: `T` is the input type and `U` is the output type.

Its constructor `NewMapper` takes `U` and `T` sequentially, which is reversed from `Mapper[T, U]`. This design assumes that Go automatically infers type `T` for the input value in constructor.

Mapper is typically capable of doing `Map` and `FlatMap`. Besides, it allows all other operations in `Reducer` for reducer is a `Mapper` without type `U`.

### Interchanging with Reducer

- Mapper to reducer: Call `mapper.NewReducerFromMapper()`(method of mapper).
- Reducer to mapper: Call `NewMapperFromReducer[U](reducer)`(top-level function).

### Map

Map allows an operation on each of the element in a `Mapper` array. It accepts a function that has 3 parameters: `item T`, `index int` and `array []T`.

Note that the function takes output value `Optional[U]`. If an empty optional is specified(see `Nothing` constructor in `Optional`), the result is not collected. Otherwise, the output for the function is collected as an array `[]U`.

It returns a new `Reducer[U]` array for further usage.

This is a mix of traditional `Map` and `FilterMap`.

- Example:
  ```go
  func main() {
  	arr := []int{1, 2, 3}
  	newArr := NewMapper[int](arr). // The int here specifies the output type U
  		Map(func (i int, _ int, _ int[]) Optional[int] {
  			if i == 1 {
  				return *Nothing[int]() // 1 in arr is skipped for return is empty
  			}
  			return *Just(i + 1) // otherwise the element is increased by 1
  		}).array
  	fmt.Println(array) // 3, 4
  }
  ```

### FlatMap

FlatMap is a flattened operation applied on Map. To be more strict, it is identity to `Map` except that the return value of the function has to be an array `[]U`.

Extra flattening is done on each of the array the function returns, by pushing the item of the array one by one into the result, thus the result being `[]U`.

It returns a new `Reducer[U]` array for further usage.

- Example:
  ```go
  func main() {
  	arr := []int{1, 2, 3}
  	newArr := NewMapper[int](arr).
  		FlatMap(func (i int, _ int, _ int[]) int[] {
  			return []int{i, i + 2} // e.g. i=1, returns [1, 3]
  		}).array
  	fmt.Println(array) // 1, 3, 2, 4, 3, 5
  }
  ```

## Reducer

`Reducer[T]` takes a generic type that specifies the input type.

Its constructor `NewReducer` takes an array. In most time, it is not needed to specify the `T` for the type inference of the given array.

Reducers are capable of doing `Filter`, `Reduce`, `ForEach`, `Slice` and `Splice`.

### Reduce

Reduce applies a function that memorizes the previous result, provides a current result and returns a future result to each of the element in an array.

It returns a single `T` value.

- Example:
  ```go
  func main() {
  	arr := []int{1, 2, 3}
  	result := NewReducer(arr).
  		Reduce(func (prev int, curr int, _ int, _ int[]) int { // previous, current, index, array
  			return prev + curr // this means doing a sum over the array
  		})
  	fmt.Println(result) // 6
  }
  ```

### Filter

Filter applies a filtering function to the array. If the function returns true, current element is kept. Otherwise, it is discarded.

It returns a filtered `Reducer[T]` for further usage.

- Example:

  ```go
  func main() {
  	arr := []int{1, 2, 3}
  	result := NewReducer(arr).
  		Filter(func (i int, _ int, _ int[]) bool {
  			if (i == 1) {
  				return false
  			}
  			return true
  		}).array
  	fmt.Println(result) // 2, 3
  }
  ```

### ForEach

Applies a generic function to each of the element in the array, regardless of return value.

It does not have a return value.

- Example:
  ```go
  func main() {
  	arr := []int{1, 2, 3}
  	newArr := make([]int, 0)
  	NewReducer(arr).
  		ForEach(func (i int, _ int, _ int[]) {
  			newArr = append(newArr, i) // an implementation of a shallow copy
  		})
  	fmt.Println(newArr) // 1, 2, 3
  }
  ```


### Slice

Borrowed from JavaScript. See [Slice](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/slice) in mdn. Different from inbuilt slices in Go, it supports negative index and returns an empty array if start and end do not overlap.

It does not change the original reducer, but returns a new sliced `Reducer[T]` for further usage.

- Example:
  ```go
  func main() {
  	arr := []int{1, 2, 3}
  	result := NewReducer(arr).
  		Slice(-2, -1).array
  	fmt.Println(result) // 2
  }
  ```

### Splice

Borrowed from JavaScript. See [Splice](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/splice) in mdn. It does 3 things:

1. Delete `deleteCount` of elements starting from `start`. `start` supports negative indexing. If `deleteCount` is too many for the array, it deletes to the end.
2. Inserts `...items` into `start`.
3. Returns the deleted array `Reducer[T]` for further usage.

Note that different from slice, this operation does deletion in place, and returns the deleted values.

- Example:
  ```go
  func main() {
  	arr := []int{1, 2, 3}
  	result := NewReducer(arr).
  		Splice(-2, 1, 4, 5).array
  	fmt.Println(result) // 2
  	fmt.Println(arr.array) // 1, 4, 5, 3
  }
  ```


## Promise

Borrowed from JavaScript and supports chaining & await. See [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) in mdn. It is actually a lower level of goroutines. It supports `Promise.all`, `Promise.any`, `NewPromise` constructor, `Then` and `Error` chainable predicates, and `Await`.

### General intro

`NewPromise` takes a concrete task `func() (T, error)` and returns a scheduled task `Promise[T]`. The scheduled task starts running immediately after calling this constructor, but is a goroutine.

`Await`, or `Promise.Await`, blocks the main goroutine and waits for a promise to complete. It is not possible to recall or chain a `Promise` further when `Await` is called on a certain instance.

`Then` chains a `Promise` to specify how the `Promise` should do next after finishing its initial job.

`AwaitAll` waits for a certain group of promises to complete, and returns a `Promise[[]T]` representing results of the promises. If an error is encountered, it returns nothing but the error.

`AwaitAny` waits for any of a promise in a certain group of promises to complete, and returns a `Promise[T]` representing the first fulfilled promise. It fails only if all promises fail.

### Designation of constructor

The constructor only accepts functions with type `func () (T, error)`. This ensures type safety for avoiding "any type" functions. To pass a parameterized function into the constructor, a "metafunction" can be specified.

- An example of specifying how many seconds to sleep:
  ```go
  sl := func(i int) func() (int, error) {
  		return func() (int, error) {
  			time.Sleep(time.Duration(i) * time.Second)
  			fmt.Println(i)
  			return i, nil
  		}
  	}
  ```

### Some examples

Using the `sl` function above, the timeline of execution of functions can be observed.

- A simple example of constructor and `Await`:

  ```go
  NewPromise(sl(1)).Await() // prints 1 after 1 sec
  ```
- An example of `AwaitAll`:

  ```go
  tasks := make([]*Promise[int], 0)
  for i := 0; i < 5; i++ {
  	tasks = append(tasks, NewPromise(sl(i + 1)))
  }
  AwaitAll(tasks...).Await() // after 1 second, each second prints 1, 2, 3 and 4
  ```
- An example of `AwaitAny`:

  ```go
  tasks := make([]*Promise[int], 0)
  for i := 4; i >= 0; i-- {
  	tasks = append(tasks, NewPromise(sl(i + 1)))
  }
  AwaitAny(tasks...).Await() // after 1 second, prints 1
  ```
- An advanced example for `then`:

  ```go
  p1 := NewPromise(sl(1)).Then(func(i int) (int, error) {
  	return sl(i + 1)() // prints 1 after 1 second, and 2 after 2 seconds then
  })
  p2 := NewPromise(sl(2)).Then(func(i int) (int, error) {
  	return sl(i + 1)() // prints 2 after 2 seconds, and 3 after 3 seconds then
  })
  AwaitAll(p1, p2).Await() // 1 second after running the program: prints 1(p1), 2: 2(p2), 3: 2(p1), 4: -, 5: 3(p2)
  ```

## Optional

Specifies a value of type `T` is empty or not. If normal function of this struct is expected, never try to directly access its fields! **Access its fields only by accessing its methods.**

The current design ensures that the value field cannot be updated once the instance is created.

### Constructors

- `Just` specifies a nonempty value to the Optional. The type is usually inferred from the value passed.
- `Nothing` specifies an empty value. The type needs to manually specify for what the `Optional` is originally intended to store.

### Methods

- `Then` chains a function to the value of Optional. If optional is empty, `Then` does nothing.
- `IsSet` checks whether the Optional is empty.
- `Value` returns the actual value of the Optional. Call this with care by checking existence of the value every time.

## Either

Similar to `Optional`, but the `Nothing` field could have any of the type like errors.

### Constructors

- `Left` specifies a left value for the Either.
- `Right` specifies a right value for the Either.

Both of the constructors accepts only one parameter, and will leave the counterpart empty.

### Methods

- `ThenIfLeft`: Applies a function to the left value, if left value exists.
- `ThenIfRight`: Applies a function to the right value, if it exists.
- `IsLeft`: Whether the Either instance is a left value.
- `IsRight`: Whether the Either instance is a right value.

## Set

Sets are collections of **unique elements**. For convenience, this data structure reuses the concept of uniqueness of keys in builtin maps.

### Constructor

`NewSet` constructs a new `Set` instance with a given typed array `[]T`.

### Methods

- `Add` adds an element to the set. If the element exists, nothing is done. Otherwise, it is inserted.
- `Delete` deletes an element in the set. If the element does not exist, nothing is done.
- `Keys` returns the `Set` instance in array `[]T`.
