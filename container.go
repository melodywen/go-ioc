package container

type Container struct {
	StructOfContainer
	BuildOfContainer
	ExtendOfContainer
}

func newContainer() *Container {
	structOfContainer := &Container{
		StructOfContainer: StructOfContainer{
			resolved:                       map[string]bool{},
			bindings:                       map[string]Bind{},
			instances:                      map[string]interface{}{},
			aliases:                        map[string]string{},
			abstractAliases:                map[string][]string{},
			extenders:                      map[string][]func(object interface{}, container *Container) interface{}{},
			reboundCallbacks:               map[string][]func(container *Container, instance interface{}){},
			globalBeforeResolvingCallbacks: nil,
			globalResolvingCallbacks:       nil,
			globalAfterResolvingCallbacks:  nil,
			beforeResolvingCallbacks:       nil,
			resolvingCallbacks:             nil,
			afterResolvingCallbacks:        nil,
		},
	}
	return structOfContainer
}
