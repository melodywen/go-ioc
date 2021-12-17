package container

import (
	util "github.com/melodywen/go-ioc/util"
)

// When Define a contextual binding.
func (container *Container) When(concrete interface{}) *ContextualBindingBuilder {
	var aliases []string
	for _, c := range util.ArrayWrap(concrete) {
		aliases = append(aliases, container.GetAlias(c))
	}
	return newContextualBindingBuilder(container, aliases)
}

// AddContextualBinding Add a contextual binding to the container.
func (container *Container) AddContextualBinding(concrete interface{}, abstract interface{}, implementation interface{}) {
	concreteIndex := container.AbstractToString(concrete)
	abstractIndex := container.AbstractToString(abstract)
	if _, ok := container.contextual[concreteIndex]; !ok {
		container.contextual[concreteIndex] = map[string]interface{}{}
	}
	container.contextual[concreteIndex][abstractIndex] = implementation
}

// getContextualConcrete 获取上下文进行解析  ,// todo : 如果是需要上下文的绑定
func (container *Container) getContextualConcrete(abstract string, buildStack []string) (concrete interface{}) {
	implementation := container.findInContextualBindings(abstract, buildStack)
	if implementation != nil {
		return implementation
	}

	// Next we need to see if a contextual binding might be bound under an alias of the
	// given abstract type. So, we will need to check if any aliases exist with this
	// type and then spin through them and check for contextual bindings on these.
	if _, ok := container.abstractAliases[abstract]; !ok {
		return nil
	}
	for _, alias := range container.abstractAliases[abstract] {
		implementation := container.findInContextualBindings(alias, buildStack)
		if implementation != nil {
			return implementation
		}
	}
	return nil
}

// Find the concrete binding for the given abstract in the contextual binding array.
func (container *Container) findInContextualBindings(abstract string, buildStack []string) interface{} {
	currentStack := buildStack[len(buildStack)-1]
	if _, ok := container.contextual[currentStack]; !ok {
		return nil
	}
	if implementation, ok := container.contextual[currentStack][abstract]; ok {
		return implementation
	}
	return nil
}
