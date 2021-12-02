package container

type ContextualBindingBuilder struct {
	container *Container // container
	concrete  []string   // contextual
	needs     string     // need obj of abstract
}

func newContextualBindingBuilder(container *Container, concrete []string) *ContextualBindingBuilder {
	return &ContextualBindingBuilder{container: container, concrete: concrete}
}

// Need Define the abstract target that depends on the context.
func (ctx *ContextualBindingBuilder) Need(abstract interface{}) *ContextualBindingBuilder {
	ctx.needs = ctx.container.AbstractToString(abstract)
	return ctx
}

// Give Define the implementation for the contextual binding
func (ctx *ContextualBindingBuilder) Give(implementation interface{}) {
	for _, concrete := range ctx.concrete {
		ctx.container.AddContextualBinding(concrete, ctx.needs, implementation)
	}
}
