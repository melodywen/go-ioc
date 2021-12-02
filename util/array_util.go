package util

import (
	"reflect"
)

// ArrayWrap 把非切片的值包裹为切片
func ArrayWrap(value interface{}) (response []interface{}) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		v := reflect.ValueOf(value)
		for i := 0; i < v.Len(); i++ {
			response = append(response, v.Index(i).Interface())
		}
	default:
		response = append(response, value)
	}
	return response
}
