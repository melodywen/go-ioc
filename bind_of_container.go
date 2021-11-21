package container

import (
	"reflect"
)

// Bind 绑定接口
func (container *Container) Bind(abstract interface{}, concrete interface{}, shared bool) {
	// 获取对应的 map key
	index := container.AbstractToString(abstract)

	// 删除老旧的实例
	container.dropStaleInstances(index)

	// If no concrete type was given, we will simply set the concrete type to the
	// abstract type. After that, the concrete type to be registered as shared
	// without being forced to state their classes in both of the parameters.
	if concrete == nil {
		concrete = abstract
	}

	// If the factory is not a Closure, it means it is just a class name which is
	// bound into this container to the abstract type and we will just wrap it
	// up inside its own Closure to give us more convenience when extending.
	if reflect.TypeOf(concrete).Kind() != reflect.Func {
		concrete = container.getClosure(index, container.AbstractToString(concrete))
	}

	// 直接进行绑定
	container.bindings[index] = Bind{shared: shared, concrete: concrete}

	// 如果是之前已经绑定过则再次重新绑定
	// If the abstract type was already resolved in this container we'll fire the
	// rebound listener so that any objects which have already gotten resolved
	// can have their copy of the object updated via the listener callbacks.
	if container.Resolved(abstract) {
		container.rebound(index)
	}
}

// dropStaleInstances
// 移除已经缓存的实例 和别名
func (container *Container) dropStaleInstances(abstract string) (ok bool) {
	if _, exist := container.instances[abstract]; exist {
		delete(container.instances, abstract)
		ok = true
	}
	if _, exist := container.aliases[abstract]; exist {
		delete(container.aliases, abstract)
		ok = true
	}
	return ok
}

// Resolved 是否已经实例化过
func (container *Container) Resolved(abstract interface{}) (ok bool) {
	if container.IsAlias(abstract) {
		abstract = container.GetAlias(abstract)
	}
	index := container.AbstractToString(abstract)

	if _, exist := container.resolved[index]; exist {
		return true
	}
	if _, exist := container.instances[index]; exist {
		return true
	}
	return false
}

// Bound Determine if the given abstract type has been bound.
func (container *Container) Bound(abstract interface{}) (ok bool) {
	index := container.AbstractToString(abstract)

	if _, exist := container.bindings[index]; exist {
		return true
	}
	if _, exist := container.instances[index]; exist {
		return true
	}
	if container.IsAlias(index) {
		return true
	}
	return false
}

// Instance Register an existing instance as shared in the container.
func (container *Container) Instance(abstract interface{}, instance interface{}) interface{} {
	index := container.AbstractToString(abstract)
	container.removeAbstractAlias(index)
	if _, ok := container.aliases[index]; ok {
		delete(container.aliases, index)
	}

	// We'll check to determine if this type has been bound before, and if it has
	// we will fire the rebound callbacks registered with the container and it
	// can be updated with consuming classes that have gotten resolved here.
	container.instances[index] = instance

	isBound := container.Bound(index)
	if isBound {
		container.rebound(index)
	}
	return instance
}
