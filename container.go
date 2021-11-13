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
			aliases:                        nil,
			abstractAliases:                nil,
			extenders:                      nil,
			reboundCallbacks:               nil,
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
