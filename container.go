package container

type Container struct {
	StructOfContainer
	BuildOfContainer
	ExtendOfContainer
}

func newContainer() *Container {
	structOfContainer := &Container{
		StructOfContainer: StructOfContainer{
			resolved:         map[string]bool{},
			bindings:         map[string]Bind{},
			instances:        map[string]interface{}{},
			aliases:          map[string]string{},
			abstractAliases:  map[string][]string{},
			extenders:        map[string][]func(object interface{}, container *Container) interface{}{},
			reboundCallbacks: map[string][]func(container *Container, instance interface{}){},

			globalBeforeResolvingCallbacks: []func(string, []interface{}, *Container){},
			globalResolvingCallbacks:       nil,
			globalAfterResolvingCallbacks:  nil,

			beforeResolvingCallbacks: map[string][]func(string, []interface{}, *Container){},
			resolvingCallbacks:       nil,
			afterResolvingCallbacks:  nil,
		},
	}
	return structOfContainer
}
