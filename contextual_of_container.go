package container

import (
	"cjw.com/melodywen/go-ioc/util"
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
func (container *Container) getContextualConcrete(abstract string) (concrete interface{}) {
	return nil
}
