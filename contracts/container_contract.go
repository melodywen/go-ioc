package contracts

import container "cjw.com/melodywen/go-ioc"

type ContainerContract interface {
	//AbstractToString 通过 abstract 进行 字符串得到 作为map key
	AbstractToString(abstract interface{}) string
	// Bound Determine if the given abstract type has been bound.
	Bound(abstract interface{}) bool
	// Alias a type to a different name.
	Alias(abstract interface{}, alias interface{})
	// Bind Register a binding with the container.
	Bind(abstract interface{}, concrete interface{}, shared bool)
	// BindIf Register a binding if it hasn't already been registered.
	BindIf(abstract interface{}, concrete interface{}, shared bool)
	// Singleton Register a binding if it hasn't already been registered.
	Singleton(abstract interface{}, concrete interface{})
	// SingletonIf Register a shared binding if it hasn't already been registered.
	SingletonIf(abstract interface{}, concrete interface{})
	// Extend "Extend" an abstract type in the container
	Extend(abstract interface{}, closure func(object interface{}, container *container.Container) interface{})
	// Instance Register an existing instance as shared in the container.
	Instance(abstract interface{}, instance interface{}) interface{}
	// AddContextualBinding Add a contextual binding to the container.
	AddContextualBinding(concrete interface{}, abstract interface{}, implementation interface{})
	// When Define a contextual binding.
	When(concrete interface{}) *container.ContextualBindingBuilder
	// Factory Get a closure to resolve the given type from the container.
	Factory(abstract interface{}) func() interface{}
	//Flush the container of all bindings and resolved instances.
	Flush()
	// Make Resolve the given type from the container.
	Make(abstract interface{}) interface{}
	// Resolved Determine if the given abstract type has been resolved.
	Resolved(abstract interface{})  bool
	// Resolving Register a new resolving callback.
	Resolving(abstract interface{}, callback func(interface{}, *container.Container))
	// AfterResolving Register a new after resolving callback for all types.
	AfterResolving(abstract interface{}, callback func(interface{}, *container.Container))
}
