package container

import "reflect"

// check concert parameters and fill
func (container *Container) getParameters(currentParams []reflect.Value, concrete interface{}, buildStack []string) []reflect.Value {
	concreteType := reflect.ValueOf(concrete).Type()
	numIn := concreteType.NumIn()
	if numIn == len(currentParams) {
		return currentParams
	}
	// if has different parameters ，will resolve params from container
	for i := 0; i < numIn; i++ {
		currentParams = append(currentParams, reflect.ValueOf(nil))
		param := concreteType.In(i)

		// 进行类型比对
		if currentParams[i].IsValid() && param == reflect.TypeOf(currentParams[i].Interface()) {
			continue
		}
		resolveParam := container.resolveDependencies(param, buildStack)
		// 在对应的位置插入实参
		currentParams = append(currentParams[:i], append([]reflect.Value{resolveParam}, currentParams[i:]...)...)
	}

	currentParams = currentParams[:numIn]
	return currentParams
}

// Resolve all the dependencies from the ReflectionParameters.
func (container *Container) resolveDependencies(parameterType reflect.Type, buildStack []string) reflect.Value {

	// 对类型进行判定
	var pkgInfo, elemInfo string
	switch parameterType.Kind() {
	//case reflect.String:
	//	fallthrough
	//case reflect.Int:
	//	fallthrough
	//case reflect.Func:
	//	panic("can not auto load param,because param type can not suppose ,please connect admin:" + parameterType.Kind().String())
	case reflect.Struct:
		pkgInfo = parameterType.PkgPath() + "." + parameterType.Name()
	case reflect.Ptr:
		elemInfo = parameterType.Elem().PkgPath() + "." + parameterType.Elem().Name()
		pkgInfo = "*" + elemInfo
	default:
		panic("can not auto load param,because param type can not suppose ,please connect admin:" + parameterType.Kind().String())
	}
	var object interface{}
	if container.Bound(pkgInfo) {
		object = container.makeWithBuildStack(container.AbstractToString(pkgInfo), nil, buildStack)

	} else if container.Bound(elemInfo) {
		object = container.makeWithBuildStack(container.AbstractToString(elemInfo), nil, buildStack)
	} else {
		panic("can not auto load param，because this interface is not registered,please connect admin")
	}
	// if need struct param ，but object is pre,so to transform data type
	if parameterType.Kind() == reflect.Struct && reflect.TypeOf(object).Kind() == reflect.Ptr {
		return reflect.ValueOf(object).Elem()
	}
	return reflect.ValueOf(object)
}

