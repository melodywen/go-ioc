package container

import (
	"reflect"
	"runtime"
)

// Factory
//  @Description: 获取一个闭包来解析容器中的给定类型。
//  @receiver container
//  @param abstract
//  @return func() any
func (container *Container) Factory(abstract any) func() any {
	return func() any {
		index := AbstractToString(abstract)
		return container.makeWithBuildStack(index, []any{}, newContainerStack())
	}
}

// Make
//  @Description: 从容器中解析给定的类型
//  @receiver container
//  @param abstract
//  @return any
func (container *Container) Make(abstract any) any {
	identifier := AbstractToString(abstract)
	return container.makeWithBuildStack(identifier, nil, newContainerStack())
}

// MakeWithParams
//  @Description: 从容器中解析给定的类型并且携带参数
//  @receiver container
//  @param abstract
//  @param parameters
//  @return any
func (container *Container) MakeWithParams(abstract any, parameters []any) any {
	identifier := AbstractToString(abstract)
	return container.makeWithBuildStack(identifier, parameters, newContainerStack())
}

// makeWithBuildStack
//  @Description: 构建上下文栈
//  @receiver container
//  @param identifier
//  @param parameters
//  @param stack
//  @return any
func (container *Container) makeWithBuildStack(identifier string, parameters []any, stack *containerStack) any {
	if stack.Stack == nil {
		pc, _, _, _ := runtime.Caller(2)
		stack.Stack = []string{runtime.FuncForPC(pc).Name()}
	}
	return container.resolve(identifier, parameters, true, stack)
}

// resolve
//  @Description: 进行解析
//  @receiver container
//  @param abstract
//  @param parameters
//  @param raiseEvents
//  @param stack
//  @return object
func (container *Container) resolve(identifier string, parameters []any, raiseEvents bool, stack *containerStack) (object any) {

	// 回调到子类
	if container.child != nil {
		container.child.ResolveCallback(identifier)
	}

	identifier = container.GetAlias(identifier)

	//首先，我们将触发任何处理“before”解析的事件处理程序特定类型。
	//这使一些钩子有机会添加各种扩展调用来改变他们。
	if raiseEvents {
		container.fireBeforeResolvingCallbacks(identifier, parameters)
	}

	concrete := container.getContextualConcrete(identifier, stack)

	needsContextualBuild := len(parameters) != 0 || concrete != nil

	//如果该类型的实例当前被管理为单例，我们将返回一个已存在的实例，而不是实例化新的实例
	//这样开发者每次都可以使用相同的对象实例。
	if _, ok := container.instances[identifier]; ok && !needsContextualBuild {
		return container.instances[identifier]
	}

	if concrete == nil {
		concrete = container.getConcrete(identifier)
	}

	//我们准备实例化注册的具体类型的实例绑定。这将实例化类型，并解析任何它的“嵌套”依赖递归地直到全部解决。
	if container.isBuildable(identifier, concrete) {
		object = container.Build(concrete, parameters, stack)
	} else {
		object = container.makeWithBuildStack(AbstractToString(concrete), parameters, stack)
	}

	//如果为该类型定义了任何扩展程序，则需要遍历它们 将它们应用到正在构建的对象。这允许扩展服务，
	//例如更改配置或装饰对象。
	for _, extender := range container.getExtenders(identifier) {
		object = extender(object, container)
	}

	//如果请求的类型注册为单例类型，我们需要缓存在“内存”中的实例，
	//这样我们可以稍后返回它，而不用创建一个一个对象的新实例。
	if container.IsShared(identifier) && !needsContextualBuild {
		container.instances[identifier] = object
	}

	if raiseEvents {
		container.fireResolvingCallbacks(identifier, object)
	}
	container.resolved[identifier] = true
	return object
}

// getConcrete
//  @Description: 获取给定抽象的具体类型。
//  @receiver container
//  @param abstract
//  @return concrete
func (container *Container) getConcrete(identifier string) (concrete any) {
	// 如果 设置了绑定的内容则返回绑定的内容
	if value, ok := container.bindings[identifier]; ok {
		return value.concrete
	}
	return identifier
}

// isBuildable 判断可以进行构建
// 1. 如果是自己等于自己则直接输出
// 2. 如果是回调函数则保留
// 3. 如果是其他的类型，说明还不够，需要进行递归获取
func (container *Container) isBuildable(identifier string, concrete any) bool {
	if reflect.DeepEqual(identifier, concrete) {
		return true
	}
	return reflect.TypeOf(concrete).Kind() == reflect.Func
}

// IsShared
//  @Description:确定给定的类型是否是共享的。
//  @receiver container
//  @param identifier
//  @return bool
func (container *Container) IsShared(identifier string) bool {
	if _, ok := container.instances[identifier]; ok {
		return true
	}
	if value, ok := container.bindings[identifier]; ok {
		return value.shared
	}
	return false
}
