package container

// Rebinding Bind a new callback to an abstract's rebind event.
func (container *Container) Rebinding(abstract interface{}, callback func(container *Container, instance interface{})) interface{} {
	index := container.AbstractToString(abstract)
	if _, ok := container.reboundCallbacks[index]; !ok {
		container.reboundCallbacks[index] = []func(container *Container, instance interface{}){}
	}
	container.reboundCallbacks[index] = append(container.reboundCallbacks[index], callback)
	if container.Bound(index) {
		return container.Make(abstract)
	}
	return nil
}

// Fire the "rebound" callbacks for the given abstract type.
func (container *Container) rebound(abstract string) {
	instance := container.Make(abstract)
	for _, callback := range container.getReboundCallbacks(abstract) {
		callback(container, instance)
	}
}

// Get the rebound callbacks for a given type.
func (container *Container) getReboundCallbacks(abstract string) []func(container *Container, instance interface{}) {
	if _, ok := container.reboundCallbacks[abstract]; ok {
		return container.reboundCallbacks[abstract]
	}
	return []func(container *Container, instance interface{}){}
}

// Get the extender callbacks for a given type.
func (container *Container) getExtenders(abstract string) []func(object interface{}, container *Container) interface{} {
	if _, ok := container.extenders[abstract]; ok {
		return container.extenders[abstract]
	}
	return []func(object interface{}, container *Container) interface{}{}
}

// ForgetExtenders Remove all the extender callbacks for a given type.
func (container *Container) ForgetExtenders(abstract interface{}) (ok bool) {
	index := container.GetAlias(abstract)
	if _, ok := container.extenders[index]; ok {
		delete(container.extenders, index)
		return true
	}
	return false
}

// Extend "Extend" an abstract type in the container
func (container *Container) Extend(abstract interface{}, closure func(object interface{}, container *Container) interface{}) {
	index := container.GetAlias(abstract)
	if _, ok := container.instances[index]; ok {
		container.instances[index] = closure(container.instances[index], container)
		container.rebound(index)
	} else {
		if _, ok := container.extenders[index]; !ok {
			container.extenders[index] = []func(object interface{}, container *Container) interface{}{}
		}
		container.extenders[index] = append(container.extenders[index], closure)
		if container.Resolved(index) {
			container.rebound(index)
		}
	}
}

//  Get the Closure to be used when building a type.
func (container *Container) getClosure(abstract string, concrete string) func() interface{} {
	return func() interface{} {
		// todo :这里有bug stack 必须向下传递
		var buildStack = []string{}
		if abstract == concrete {
			return container.Build(concrete, []interface{}{}, buildStack)
		}


		return container.resolve(concrete, []interface{}{}, false, buildStack)
	}
}

//Fire an array of callbacks with an object.
func (container *Container) fireBeforeCallbackArray(
	abstract string,
	parameters []interface{},
	callbacks []func(string, []interface{}, *Container)) {
	for _, callback := range callbacks {
		callback(abstract, parameters, container)
	}
}

// Fire all the before resolving callbacks.
func (container *Container) fireBeforeResolvingCallbacks(abstract string, parameters []interface{}) {
	container.fireBeforeCallbackArray(abstract, parameters, container.globalBeforeResolvingCallbacks)

	for index, callbacks := range container.beforeResolvingCallbacks {
		// todo 可以做类似spring 的aop 操作
		if index == abstract {
			container.fireBeforeCallbackArray(abstract, parameters, callbacks)
		}
	}
}

// BeforeResolving Register a new before resolving callback for all types.
func (container *Container) BeforeResolving(abstract interface{}, callback func(string, []interface{}, *Container)) {
	if abstract == nil {
		container.globalBeforeResolvingCallbacks = append(container.globalBeforeResolvingCallbacks, callback)
		return
	}
	index := container.GetAlias(abstract)

	if _, ok := container.beforeResolvingCallbacks[index]; !ok {
		container.beforeResolvingCallbacks[index] = []func(string, []interface{}, *Container){}
	}
	container.beforeResolvingCallbacks[index] = append(container.beforeResolvingCallbacks[index], callback)
}

// Fire an array of callbacks with an object.
func (container *Container) fireCallbackArray(object interface{}, callbacks []func(interface{}, *Container)) {
	for _, callback := range callbacks {
		callback(object, container)
	}
}

// Get all callbacks for a given type.
func (container *Container) getCallbacksForType(abstract string, object interface{}, callbacks map[string][]func(interface{}, *Container)) []func(interface{}, *Container) {
	var results []func(interface{}, *Container)
	for index, callback := range callbacks {
		// todo 可以做类似spring 的aop 操作
		if index == abstract {
			results = append(results, callback...)
		}
	}
	return results
}

// Resolving Register a new resolving callback.
func (container *Container) Resolving(abstract interface{}, callback func(interface{}, *Container)) {
	if abstract == nil {
		container.globalResolvingCallbacks = append(container.globalResolvingCallbacks, callback)
		return
	}
	index := container.GetAlias(abstract)
	if _, ok := container.resolvingCallbacks[index]; !ok {
		container.resolvingCallbacks[index] = []func(interface{}, *Container){}
	}
	container.resolvingCallbacks[index] = append(container.resolvingCallbacks[index], callback)
}

// AfterResolving Register a new after resolving callback for all types.
func (container *Container) AfterResolving(abstract interface{}, callback func(interface{}, *Container)) {
	if abstract == nil {
		container.globalAfterResolvingCallbacks = append(container.globalAfterResolvingCallbacks, callback)
		return
	}
	index := container.GetAlias(abstract)
	if _, ok := container.afterResolvingCallbacks[index]; !ok {
		container.afterResolvingCallbacks[index] = []func(interface{}, *Container){}
	}
	container.afterResolvingCallbacks[index] = append(container.afterResolvingCallbacks[index], callback)
}

// Fire all the resolving callbacks.
func (container *Container) fireResolvingCallbacks(abstract string, object interface{}) {
	container.fireCallbackArray(object, container.globalResolvingCallbacks)

	container.fireCallbackArray(
		object,
		container.getCallbacksForType(abstract, object, container.resolvingCallbacks),
	)
	container.fireAfterResolvingCallbacks(abstract, object)
}

// Fire all the resolving callbacks.
func (container *Container) fireAfterResolvingCallbacks(abstract string, object interface{}) {
	container.fireCallbackArray(object, container.globalAfterResolvingCallbacks)
	container.fireCallbackArray(
		object,
		container.getCallbacksForType(abstract, object, container.afterResolvingCallbacks),
	)
}
