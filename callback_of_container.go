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
func (container *Container) getExtenders(abstract string) []interface{} {
	if _, ok := container.extenders[abstract]; ok {
		return container.extenders[abstract]
	}
	return []interface{}{}
}

// Remove all of the extender callbacks for a given type.
func (container *Container) forgetExtenders(abstract string) (ok bool) {
	if _, ok := container.extenders[abstract]; ok {
		delete(container.extenders, abstract)
		return true
	}
	return false
}

// Extend "Extend" an abstract type in the container
func (container *Container) Extend(abstract string, closure func(interface{}, *Container) interface{}) (ok bool) {
	abstract = container.GetAlias(abstract)

	// TODO 待完成
	if _, ok := container.instances[abstract]; ok {
		container.instances[abstract] = closure(container.instances[abstract], container)
	}

	return false
}
