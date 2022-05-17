package container

// ContextualBindingBuilder
//  @Description:
type ContextualBindingBuilder struct {
	container *Container // container
	concrete  []string   // contextual
	needs     string     // need obj of abstract
}

// newContextualBindingBuilder
//  @Description: 创建一个新的上下文绑定构建器
//  @param container
//  @param concrete
//  @return *ContextualBindingBuilder
func newContextualBindingBuilder(container *Container, concrete []string) *ContextualBindingBuilder {
	return &ContextualBindingBuilder{container: container, concrete: concrete}
}

// Need
//  @Description: 定义依赖于上下文的抽象目标
//  @receiver ctx
//  @param abstract
//  @return *ContextualBindingBuilder
func (ctx *ContextualBindingBuilder) Need(abstract interface{}) *ContextualBindingBuilder {
	ctx.needs = AbstractToString(abstract)
	return ctx
}

// Give
//  @Description: 定义上下文绑定的实现
//  @receiver ctx
//  @param implementation
func (ctx *ContextualBindingBuilder) Give(implementation interface{}) {
	for _, concrete := range ctx.concrete {
		ctx.container.AddContextualBinding(concrete, ctx.needs, implementation)
	}
}

// When
//  @Description: 定义上下文绑定
//  @receiver container
//  @param concretes
//  @return *ContextualBindingBuilder
func (container *Container) When(concretes []any) *ContextualBindingBuilder {
	var aliases []string
	for _, concrete := range concretes {
		aliases = append(aliases, container.GetAlias(concrete))
	}
	return newContextualBindingBuilder(container, aliases)
}

// AddContextualBinding
//  @Description: 向容器添加上下文绑定
//  @receiver container
//  @param concrete
//  @param abstract
//  @param implementation
func (container *Container) AddContextualBinding(concrete, abstract, implementation any) {
	concreteIndex := AbstractToString(concrete)
	abstractIndex := AbstractToString(abstract)
	if _, ok := container.contextual[concreteIndex]; !ok {
		container.contextual[concreteIndex] = map[string]interface{}{}
	}
	container.contextual[concreteIndex][abstractIndex] = implementation
}

// getContextualConcrete
//  @Description: 获取给定抽象的上下文具体绑定
//  @receiver container
//  @param abstract
//  @param buildStack
//  @return concrete
func (container *Container) getContextualConcrete(abstract string, stack *containerStack) (concrete interface{}) {
	implementation := container.findInContextualBindings(abstract, stack)
	if implementation != nil {
		return implementation
	}

	//接下来我们需要看看上下文绑定是否可以在给定的抽象类型。
	//所以，我们需要检查是否有化名存在类型，然后遍历它们并检查它们的上下文绑定。
	if _, ok := container.abstractAliases[abstract]; !ok {
		return nil
	}
	for _, alias := range container.abstractAliases[abstract] {
		implementation := container.findInContextualBindings(alias, stack)
		if implementation != nil {
			return implementation
		}
	}
	return nil
}

// findInContextualBindings
//  @Description: 在上下文绑定数组中查找给定抽象的具体绑定。
//  @receiver container
//  @param abstract
//  @param stack
//  @return interface{}
func (container *Container) findInContextualBindings(abstract string, stack *containerStack) interface{} {
	currentStack := stack.Stack[len(stack.Stack)-1]
	if _, ok := container.contextual[currentStack]; !ok {
		return nil
	}
	if implementation, ok := container.contextual[currentStack][abstract]; ok {
		return implementation
	}
	return nil
}
