package contracts

import "github.com/melodywen/go-ioc/container"

// ContainerContract
// @Description: 容器接口
type ContainerContract interface {
	Bound(abstract any) (ok bool)
	Alias(abstract, alias any)
	Bind(abstract, concrete any, shared bool)
	BindIf(abstract, concrete any, shared bool)
	Singleton(abstract, instance any)
	SingletonIf(abstract, instance any)
	Extend(abstract any, closure func(object any, container *container.Container) any)
	Instance(abstract, instance any) any
	AddContextualBinding(concrete, abstract, implementation any)
	When(concretes []any) *container.ContextualBindingBuilder
	Factory(abstract any) func() any
	Flush()
	Make(abstract any) any
	MakeWithParams(abstract any, parameters []any) any
	Resolved(abstract any) (ok bool)
	BeforeResolving(abstract any, callback func(string, []any, *container.Container))
	Resolving(abstract any, callback func(any, *container.Container))
	AfterResolving(abstract any, callback func(any, *container.Container))
}
