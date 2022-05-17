package container

import (
	"github.com/melodywen/supports/exceptions"
	"reflect"
)

// BindIf
//  @Description:注册尚未注册的绑定
//  @receiver container
//  @param abstract
//  @param concrete
//  @param shared
func (container *Container) BindIf(abstract, concrete any, shared bool) {
	if !container.Bound(abstract) {
		container.Bind(abstract, concrete, shared)
	}
}

// Bind
//  @Description:向容器注册一个绑定
//  @receiver container
//  @param abstract
//  @param concrete
//  @param shared
func (container *Container) Bind(abstract any, concrete any, shared bool) {
	container.storeMutex.Lock()
	container.aliasesMutex.Lock()

	identifier := AbstractToString(abstract)
	container.dropStaleInstances(identifier)

	//如果没有给出具体类型，我们将简单地将具体类型设置为抽象类型。
	//之后，要注册为共享的具体类型
	//而不是强制在两个形参中声明它们的类。
	if concrete == nil {
		concrete = abstract
	}

	//如果工厂不是一个闭包，这意味着它只是一个类名
	//绑定到这个容器的抽象类型，我们只会包装它
	//在它自己的闭包中给我们更多的方便时扩展。
	if reflect.TypeOf(concrete).Kind() != reflect.Func {
		concrete = container.getClosure(identifier, AbstractToString(concrete))
	}

	container.bindings[identifier] = bind{shared: shared, concrete: concrete}

	container.storeMutex.Unlock()
	container.aliasesMutex.Unlock()

	//如果抽象类型已经在容器中解析，则触发
	//反弹监听器，以便任何对象已经得到解决
	//可以通过监听器回调来更新它们的对象副本
	if container.Resolved(abstract) {
		container.rebound(identifier)
	}
}

// Resolved
//  @Description: 确定给定的抽象类型是否已解析
//  @receiver container
//  @param abstract
//  @return ok
func (container *Container) Resolved(abstract any) (ok bool) {
	identifier := AbstractToString(abstract)
	if container.IsAlias(identifier) {
		identifier = container.GetAlias(identifier)
	}
	if _, exist := container.resolved[identifier]; exist {
		return true
	}
	if _, exist := container.instances[identifier]; exist {
		return true
	}
	return false
}

// dropStaleInstances
//  @Description: 删除所有过时的实例和别名
//  @receiver container
//  @param abstract
func (container *Container) dropStaleInstances(identifier string) {
	delete(container.instances, identifier)
	delete(container.aliases, identifier)
}

// getClosure
//  @Description:获取在构建类型时使用的闭包
//  @receiver container
//  @param identifier
//  @param concrete
//  @return func(containerStack) any
func (container *Container) getClosure(identifier, concrete string) func(*containerStack) any {
	return func(stack *containerStack) any {
		if identifier == concrete {
			container.Build(concrete, nil, stack)
		}
		return container.resolve(concrete, []interface{}{}, false, stack)
	}
}

// IsAlias
//  @Description:是否为别名
//  @receiver container
//  @param abstract
//  @return ok
func (container *Container) IsAlias(abstract any) bool {
	_, ok := container.aliases[AbstractToString(abstract)]
	return ok
}

// GetAlias
//  @Description: 如果可用，获取摘要的别名
//  @receiver container
//  @param abstract
//  @return string
func (container *Container) GetAlias(abstract any) string {
	identifier := AbstractToString(abstract)
	if _, ok := container.aliases[identifier]; ok {
		return container.GetAlias(container.aliases[identifier])
	}
	return identifier
}

// Alias
//  @Description: 将类型别名为不同的名称
//  @receiver container
//  @param abstract
//  @param alias
func (container *Container) Alias(abstract, alias any) {
	container.aliasesMutex.Lock()
	defer func() {
		container.aliasesMutex.Unlock()
	}()

	abstractStr := AbstractToString(abstract)
	aliasStr := AbstractToString(alias)
	if aliasStr == abstractStr {
		panic(exceptions.NewInvalidParamError("[{abstract}] is aliased to itself."))
	}
	container.aliases[aliasStr] = abstractStr
	container.abstractAliases[abstractStr] = append(container.abstractAliases[abstractStr], aliasStr)
}

// removeAbstractAlias
//  @Description: 从上下文绑定别名缓存中删除别名
//  @receiver container
//  @param search
func (container *Container) removeAbstractAlias(search string) {
	if _, ok := container.aliases[search]; !ok {
		return
	}
	for abstract, aliases := range container.abstractAliases {
		for index, alias := range aliases {
			if alias == search {
				container.abstractAliases[abstract] = append(aliases[:index], aliases[index+1:]...)
			}
		}
	}
}

// Bound
//  @Description:确定给定的抽象类型是否已被绑定
//  @receiver container
//  @param abstract
//  @return ok
func (container *Container) Bound(abstract any) (ok bool) {
	identifier := AbstractToString(abstract)
	if _, exist := container.bindings[identifier]; exist {
		return true
	}
	if _, exist := container.instances[identifier]; exist {
		return true
	}
	if container.IsAlias(identifier) {
		return true
	}
	return false
}

// Instance
//  @Description: 将现有实例注册为容器中的共享实例
//  @receiver container
//  @param abstract
//  @param instance
//  @return interface{}
func (container *Container) Instance(abstract, instance any) any {
	identifier := AbstractToString(abstract)
	container.removeAbstractAlias(identifier)

	delete(container.aliases, identifier)

	isBound := container.Bound(identifier)

	//我们将检查该类型之前是否绑定过，是否已经绑定
	//我们将触发注册到容器和它的反弹回调函数
	//可以通过在这里解析的消费类来更新。
	container.instances[identifier] = instance

	if isBound {
		container.rebound(identifier)
	}
	return instance
}

// Singleton
//  @Description:在容器中注册一个共享绑定
//  @receiver container
//  @param abstract
//  @param concrete
func (container *Container) Singleton(abstract, concrete any) {
	container.Bind(abstract, concrete, true)
}

// SingletonIf
//  @Description:注册一个尚未注册的共享绑定
//  @receiver container
//  @param abstract
//  @param concrete
func (container *Container) SingletonIf(abstract, concrete any) {
	if !container.Bound(abstract) {
		container.Singleton(abstract, concrete)
	}
}

// Flush
//  @Description: 刷新所有绑定和已解析实例的容器
//  @receiver container
func (container *Container) Flush() {
	container.aliases = map[string]string{}
	container.resolved = map[string]bool{}
	container.bindings = map[string]bind{}
	container.instances = map[string]interface{}{}
	container.abstractAliases = map[string][]string{}
	container.contextual = map[string]map[string]any{}
}

// rebound
//  @Description: 为给定的抽象类型触发“反弹”回调
//  @receiver container
//  @param abstract
func (container *Container) rebound(abstract string) {
	instance := container.Make(abstract)
	for _, callback := range container.getReboundCallbacks(abstract) {
		callback(container, instance)
	}
}
