package container

import (
	"reflect"
	"runtime"
	"strconv"
)

// ExtendOfContainer container struct component
type ExtendOfContainer struct {
}

// AbstractToString 通过 abstract 进行 字符串得到 作为map key
// 可以是 int string struct、ptr
// 如果是标量则直接使用标量 ，其他的类型则是包的地址
func (extendOfContainer *ExtendOfContainer) AbstractToString(abstract interface{}) (response string) {
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
		panic("checkAbstract error")
	}
	return response
}
