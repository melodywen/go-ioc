package container

import (
	"reflect"
	"runtime"
	"strconv"
)

type CommonOfContainer struct {
}

// AbstractToString 通过 abstract 进行 字符串得到 作为map key
// 可以是 int string struct、ptr
func (common *CommonOfContainer) AbstractToString(abstract interface{}) (response string) {
	classInfo := reflect.TypeOf(abstract)
	switch classInfo.Kind() {
	case reflect.String:
		response = abstract.(string)
	case reflect.Int:
		response = strconv.Itoa(abstract.(int))
	case reflect.Struct:
		response = classInfo.PkgPath() + "." + classInfo.Name()
	case reflect.Func:
		response = runtime.FuncForPC(reflect.ValueOf(abstract).Pointer()).Name()
	case reflect.Ptr:
		response = "*" + classInfo.Elem().PkgPath() + "." + classInfo.Elem().Name()
	default:
		panic("checkAbstract error ")
	}
	return response
}
