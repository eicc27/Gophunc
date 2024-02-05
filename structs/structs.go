package structs

import (
	"reflect"

	A "github.com/eicc27/Gophunc/array"
)

// Gets the keys of a field from a struct.
// If the object is not a struct, returns an empty array.
func Keys(object any) *A.TypedArray[string, any] {
	if reflect.TypeOf(object).Kind() != reflect.Struct {
		return A.New[string]()
	}
	values := reflect.ValueOf(object)
	return A.WithType[string](A.TypedCount(values.NumField())).SimpleMap(func(t int) string {
		return values.Type().Field(t).Name
	})
}

// Gets a value from the object with given string key.
func ValueOf(object any, key string) any {
	if reflect.TypeOf(object).Kind() != reflect.Struct {
		return nil
	}
	values := reflect.ValueOf(object)
	return values.FieldByName(key).Interface()
}
