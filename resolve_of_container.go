package container

import (
	"reflect"
	"runtime"
)

// Factory Get a closure to resolve the given type from the container.
func (container *Container) Factory(abstract interface{}) func() interface{} {
	return func() interface{} {
		index := container.AbstractToString(abstract)
		return container.makeWithBuildStack(index, []interface{}{}, nil)
	}
}

// MakeWithParams Resolve the given type from the container with parameters.
func (container *Container) MakeWithParams(abstract interface{}, parameters []interface{}) interface{} {
	index := container.AbstractToString(abstract)
	return container.makeWithBuildStack(index, parameters, nil)
}

// Make Resolve the given type from the container.
func (container *Container) Make(abstract interface{}) interface{} {
	index := container.AbstractToString(abstract)
	return container.makeWithBuildStack(index, []interface{}{}, nil)
}

func (container *Container) makeWithBuildStack(abstract string, parameters []interface{}, buildStack []string) interface{} {
	index := container.AbstractToString(abstract)
	if buildStack == nil {
		pc, _, _, _ := runtime.Caller(2)
		buildStack = []string{runtime.FuncForPC(pc).Name()}
	}
	return container.resolve(index, parameters, true, buildStack)
}

// resolve 进行解析
func (container *Container) resolve(abstract string, parameters []interface{}, raiseEvents bool, buildStack []string) (object interface{}) {
	// callback child method callback method
	if container.child != nil {
		relValue := reflect.ValueOf(container.child)
		method := relValue.MethodByName("ResolveCallback")
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(abstract)})
		}
	}

	abstract = container.GetAlias(abstract)

	// First we'll fire any event handlers which handle the "before" resolving of
	// specific types. This gives some hooks the chance to add various extends
	// calls to change the resolution of objects that they're interested in.
	if raiseEvents {
		container.fireBeforeResolvingCallbacks(abstract, parameters)
	}

	concrete := container.getContextualConcrete(abstract, buildStack)

	needsContextualBuild := len(parameters) != 0 || concrete != nil

	// If an instance of the type is currently being managed as a singleton we'll
	// just return an existing instance instead of instantiating new instances
	// so the developer can keep using the same objects instance every time.
	if _, ok := container.instances[abstract]; ok && !needsContextualBuild {
		return container.instances[abstract]
	}

	if concrete == nil {
		concrete = container.getConcrete(abstract)
	}

	// We're ready to instantiate an instance of the concrete type registered for
	// the binding. This will instantiate the types, as well as resolve any of
	// its "nested" dependencies recursively until all have gotten resolved.
	if container.isBuildable(abstract, concrete) {
		object = container.Build(concrete, parameters, buildStack)
	} else {
		object = container.makeWithBuildStack(container.AbstractToString(concrete), parameters, buildStack)
	}

	// If we defined any extenders for this type, we'll need to spin through them
	// and apply them to the object being built. This allows for the extension
	// of services, such as changing configuration or decorating the object.
	for _, extender := range container.getExtenders(abstract) {
		object = extender(object, container)
	}

	// If the requested type is registered as a singleton we'll want to cache off
	// the instances in "memory" so we can return it later without creating an
	// entirely new instance of an object on each subsequent request for it.
	if container.IsShared(abstract) && !needsContextualBuild {
		container.instances[abstract] = object
	}

	if raiseEvents {
		container.fireResolvingCallbacks(abstract, object)
	}

	container.resolved[abstract] = true
	return object

}

// GetConcrete 通过接口获取具体实现
func (container *Container) getConcrete(abstract string) (concrete interface{}) {
	// 如果 设置了绑定的内容则返回绑定的内容
	if value, ok := container.bindings[abstract]; ok {
		return value.concrete
	}
	return abstract
}

// IsShared 判断这个接口是否为共享
func (container *Container) IsShared(abstract string) bool {
	if _, ok := container.instances[abstract]; ok {
		return true
	}
	if value, ok := container.bindings[abstract]; ok {
		return value.shared
	}
	return false
}
