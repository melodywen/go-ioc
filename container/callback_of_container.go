package container

// BeforeResolving
//  @Description: 在为所有类型解析回调之前注册一个新对象
//  @receiver container
//  @param abstract
//  @param callback
func (container *Container) BeforeResolving(abstract any, callback func(string, []any, *Container)) {
	container.callbackMutex.Lock()
	defer container.callbackMutex.Unlock()
	if abstract == nil {
		container.globalBeforeResolvingCallbacks = append(container.globalBeforeResolvingCallbacks, callback)
		return
	}
	identifier := container.GetAlias(abstract)
	container.beforeResolvingCallbacks[identifier] = append(container.beforeResolvingCallbacks[identifier], callback)
}

// Resolving
//  @Description: 注册一个新的解析回调。
//  @receiver container
//  @param abstract
//  @param callback
func (container *Container) Resolving(abstract any, callback func(any, *Container)) {
	container.callbackMutex.Lock()
	defer container.callbackMutex.Unlock()
	if abstract == nil {
		container.globalResolvingCallbacks = append(container.globalResolvingCallbacks, callback)
		return
	}
	identifier := container.GetAlias(abstract)
	container.resolvingCallbacks[identifier] = append(container.resolvingCallbacks[identifier], callback)
}

// AfterResolving
//  @Description: 为所有类型解析回调后注册一个新对象
//  @receiver container
//  @param abstract
//  @param callback
func (container *Container) AfterResolving(abstract any, callback func(any, *Container)) {
	container.callbackMutex.Lock()
	defer container.callbackMutex.Unlock()
	if abstract == nil {
		container.globalAfterResolvingCallbacks = append(container.globalAfterResolvingCallbacks, callback)
		return
	}
	identifier := container.GetAlias(abstract)
	container.afterResolvingCallbacks[identifier] = append(container.afterResolvingCallbacks[identifier], callback)
}

// fireBeforeCallbackArray
//  @Description:对象触发一个回调数组
//  @receiver container
//  @param identifier
//  @param parameters
//  @param callbacks
func (container *Container) fireBeforeCallbackArray(identifier string, parameters []any, callbacks []func(string, []any, *Container)) {
	for _, callback := range callbacks {
		callback(identifier, parameters, container)
	}
}

// fireBeforeResolvingCallbacks
//  @Description: 在解析回调之前触发所有的回调
//  @receiver container
//  @param identifier
//  @param parameters
func (container *Container) fireBeforeResolvingCallbacks(identifier string, parameters []any) {
	container.fireBeforeCallbackArray(identifier, parameters, container.globalBeforeResolvingCallbacks)
	for index, callbacks := range container.beforeResolvingCallbacks {
		if index == identifier {
			container.fireBeforeCallbackArray(identifier, parameters, callbacks)
			break
		}
	}
}

// fireCallbackArray
//  @Description: 用对象触发一个回调数组。
//  @receiver container
//  @param object
//  @param callbacks
func (container *Container) fireCallbackArray(object any, callbacks []func(any, *Container)) {
	for _, callback := range callbacks {
		callback(object, container)
	}
}

// fireResolvingCallbacks
//  @Description: 触发所有解析回调。
//  @receiver container
//  @param identifier
//  @param object
func (container *Container) fireResolvingCallbacks(identifier string, object any) {
	container.fireCallbackArray(object, container.globalResolvingCallbacks)
	for index, callbacks := range container.resolvingCallbacks {
		if index == identifier {
			container.fireCallbackArray(object, callbacks)
			break
		}
	}
	container.fireAfterResolvingCallbacks(identifier, object)
}

// fireAfterResolvingCallbacks
//  @Description: 触发所有的after解析回调。
//  @receiver container
//  @param identifier
//  @param object
func (container *Container) fireAfterResolvingCallbacks(identifier string, object any) {
	container.fireCallbackArray(object, container.globalAfterResolvingCallbacks)
	for index, callbacks := range container.afterResolvingCallbacks {
		if index == identifier {
			container.fireCallbackArray(object, callbacks)
			break
		}
	}
}

// getExtenders
//  @Description: 获取给定类型的扩展程序回调
//  @receiver container
//  @param identifier
//  @return []func(object
func (container *Container) getExtenders(identifier string) []func(object any, container *Container) any {
	return container.extenders[identifier]
}

// ForgetExtenders
//  @Description: Remove all of the extender callbacks for a given type.
//  @receiver container
//  @param abstract
func (container *Container) ForgetExtenders(abstract any) {
	delete(container.extenders, container.GetAlias(abstract))
}

// Extend
//  @Description: 扩展“容器中的抽象类型”
//  @receiver container
//  @param abstract
//  @param closure
func (container *Container) Extend(abstract any, closure func(object any, container *Container) any) {
	container.callbackMutex.Lock()
	defer container.callbackMutex.Unlock()
	identifier := container.GetAlias(abstract)
	if _, ok := container.instances[identifier]; ok {
		container.instances[identifier] = closure(container.instances[identifier], container)
		container.rebound(identifier)
	} else {
		container.extenders[identifier] = append(container.extenders[identifier], closure)
		if container.Resolved(identifier) {
			container.rebound(identifier)
		}
	}
}

// Rebinding
//  @Description: 将一个新的回调绑定到一个抽象的重新绑定事件。
//  @receiver container
//  @param abstract
//  @param callback
//  @return any
func (container *Container) Rebinding(abstract any, callback func(container *Container, instance any)) any {
	container.callbackMutex.Lock()
	defer container.callbackMutex.Unlock()

	identifier := container.GetAlias(abstract)
	container.reboundCallbacks[identifier] = append(container.reboundCallbacks[identifier], callback)
	if container.Bound(identifier) {
		return container.Make(abstract)
	}
	return nil
}

// getReboundCallbacks
//  @Description: 获取给定类型的反弹回调
//  @receiver container
//  @param identifier
//  @return []func(container
func (container *Container) getReboundCallbacks(identifier string) []func(container *Container, instance any) {
	return container.reboundCallbacks[identifier]
}
