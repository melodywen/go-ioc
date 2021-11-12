package container

type Container struct {
	StructOfContainer
	BuildOfContainer
	ExtendOfContainer
}

func newContainer() *Container {
	structOfContainer := &Container{
		StructOfContainer: StructOfContainer{
			resolved:                       nil,
			bindings:                       map[string]Bind{},
			instances:                      nil,
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
