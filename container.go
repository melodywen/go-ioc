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
			globalResolvingCallbacks:       []func(interface{}, *Container){},
			globalAfterResolvingCallbacks:  []func(interface{}, *Container){},

			beforeResolvingCallbacks: map[string][]func(string, []interface{}, *Container){},
			resolvingCallbacks:       map[string][]func(interface{}, *Container){},
			afterResolvingCallbacks:  map[string][]func(interface{}, *Container){},

			contextual: map[string]map[string]interface{}{},
		},
	}
	return structOfContainer
}
