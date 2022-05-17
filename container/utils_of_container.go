package container

import (
	"encoding/json"
	"github.com/melodywen/supports/exceptions"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// AbstractToString
//  @Description: 提取一个变量的标识符
//  @param abstract
//  @return identifier
func AbstractToString(abstract any) (identifier string) {
	return cachedAbstractToString(abstract)
}

var abstractToStringCache sync.Map // map[reflect.Type]string

// cachedAbstractToString
//  @Description:cache abstract string
//  @param abstract
//  @return string
func cachedAbstractToString(abstract any) string {
	t := reflect.TypeOf(abstract)
	if t.Kind() == reflect.String {
		return abstract.(string)
	}

	if f, ok := abstractToStringCache.Load(t); ok {
		return f.(string)
	}
	f, _ := abstractToStringCache.LoadOrStore(t, abstractToStringHandle(abstract, t))
	return f.(string)
}

// cachedAbstractTypeToString
//  @Description:cache abstract type string
//  @param t
//  @return string
func cachedAbstractTypeToString(t reflect.Type) string {
	if f, ok := abstractToStringCache.Load(t); ok {
		return f.(string)
	}
	f, _ := abstractToStringCache.LoadOrStore(t, abstractToStringHandle(nil, t))
	return f.(string)
}

// abstractToStringHandle
//  @Description:abstract to string handle method
//  @param abstract
//  @param t
//  @return identifier
func abstractToStringHandle(abstract any, t reflect.Type) (identifier string) {
	isAdd := false
	switch t.Kind() {
	case reflect.String:
		identifier = abstract.(string)
	case reflect.Func:
		identifier = runtime.FuncForPC(reflect.ValueOf(abstract).Pointer()).Name()
	case reflect.Ptr:
		t = t.Elem()
		isAdd = true
		fallthrough
	case reflect.Struct, reflect.Interface:
		identifier = t.PkgPath() + "." + t.Name()
	default:
		panic(exceptions.NewInvalidParamError("invalid data type"))
	}
	if isAdd {
		identifier = "*" + identifier
	}
	return identifier
}

// isConstructMethod
//  @Description: 判定是否为构造方法
//  @param concrete
//  @return bool
func isConstructMethod(identifier string) bool {
	index := strings.Split(identifier, ".")
	methodName := index[len(index)-1]
	if strings.HasPrefix(methodName, "New") || strings.HasPrefix(methodName, "new") {
		return true
	}
	return false
}

// errorOther
//  @Description: 错误信息格式化
//  @param subject
//  @return string
func errorOther(subject map[string]any) string {
	response, e := json.Marshal(subject)
	if e != nil {
		return e.Error()
	}
	return string(response)
}
