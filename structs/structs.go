package structs

import (
	"reflect"

	"github.com/eicc27/Gophunc/array"
)

// Gets the keys of a field from a struct.
// If the object is not a struct, returns an empty array.
func Keys(o any) *array.TypedArray[string, any] {
	if reflect.TypeOf(o).Kind() != reflect.Struct {
		return array.NewTypedArray[string]()
	}
	values := reflect.ValueOf(o)
	return array.WithType[string](array.TypedCount(values.NumField())).SimpleMap(func(t int) string {
		return values.Type().Field(t).Name
	})
}

// Gets a value from the object with given string key.
func ValueOf(o any, k string) any {
	if reflect.TypeOf(o).Kind() != reflect.Struct {
		return nil
	}
	values := reflect.ValueOf(o)
	return values.FieldByName(k).Interface()
}