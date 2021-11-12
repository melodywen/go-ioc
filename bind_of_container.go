package container

import "reflect"

// Bind 绑定接口
func (container *Container) Bind(abstract interface{}, concrete interface{}, shared bool) {
	// 获取对应的 map key
	index := container.AbstractToString(abstract)

	// 删除老旧的实例
	container.dropStaleInstances(index)

	if concrete == nil {
		concrete = abstract
	}

	if reflect.TypeOf(concrete).Kind() != reflect.Func {
		// todo 等待完成
		//concrete = container.getClosure(index, concrete)
	}

	// 直接进行绑定
	container.bindings[index] = Bind{shared: shared, concrete: concrete}

	// 如果是之前已经绑定过则再次重新绑定 todo 等待完成
	//if container.Resolved(abstract) {
	//	container.Rebound(index)
	//}
}

// dropStaleInstances
// 移除已经缓存的实例 和别名
func (container *Container) dropStaleInstances(abstract string) (ok bool) {
	if _, ok = container.instances[abstract]; ok {
		delete(container.instances, abstract)
	}
	if _, ok = container.aliases[abstract]; ok {
		delete(container.aliases, abstract)
	}
	return ok
}
