package container

import (
	"reflect"
)

//type Bind struct {
//	shared   bool
//	concrete interface{}
//}

type Container struct {
	instances map[string]interface{} // 绑定的实例 ， 如果他是单例模式则全部存储到这里面
	bindings  map[string]Bind        // 绑定的策略及其配置
	resolved  map[string]bool        // 是否最终解析成功
}

func NewContainer() *Container {
	instances := make(map[string]interface{})
	bindings := make(map[string]Bind)
	resolved := make(map[string]bool)
	obj := &Container{instances: instances, bindings: bindings, resolved: resolved}
	return obj
}

// Singleton 绑定一个单例
func (container *Container) Singleton(abstract interface{}) {
	container.SingletonWithConcrete(abstract, abstract)
}

// SingletonWithConcrete 绑定一个共享的接口
func (container *Container) SingletonWithConcrete(abstract interface{}, concrete interface{}) {
	container.Bind(abstract, concrete, true)
}

//  获取对应的接口名称
// abstract 目前完成了 对象 、 接口 和字符串
func (container *Container) checkAbstract(abstract interface{}) string {
	classInfo := reflect.TypeOf(abstract)
	var pkgPath string
	var name string
	switch classInfo.Kind() {
	case reflect.String:
		return abstract.(string)
	case reflect.Struct,reflect.Func:
		pkgPath = classInfo.PkgPath()
		name = classInfo.Name()
	case reflect.Ptr:
		classInfo = classInfo.Elem()
		pkgPath = classInfo.PkgPath()
		name = classInfo.Name()
	default:
		panic("checkAbstract error ")
	}
	return pkgPath + "/" + name
}

// Bind 绑定一个实例
func (container *Container) Bind(abstract interface{}, concrete interface{}, shared bool) {
	// get abstract value to string ,set to component index
	index := container.checkAbstract(abstract)

	container.DropStaleInstances(index)

	// todo: 如果他不是一个闭包函数，则把它扩展为一个闭包函数，后面扩展就方便很多

	container.bindings[index] = Bind{
		shared:   shared,
		concrete: concrete,
	}

	// 如果是之前已经绑定过则再次重新绑定,
	if container.Resolved(abstract) {
		container.rebound(index)
	}
}

// DropStaleInstances 删除一个老的实例
func (container *Container) DropStaleInstances(abstract interface{}) {
	index := container.checkAbstract(abstract)
	if _, ok := container.instances[index]; ok {
		delete(container.instances, index)
	}
}

// Make 对外暴露make 方法
func (container *Container) Make(abstract interface{}) interface{} {
	return container.resolve(abstract, []interface{}{}, true)
}

// MakeWithParams 对外暴露make 方法并且
func (container *Container) MakeWithParams(abstract interface{}, parameters []interface{}) interface{} {
	return container.resolve(abstract, parameters, true)
}

// Resolved 是否已经绑定过接口
func (container *Container) Resolved(abstract interface{}) bool {
	index := container.checkAbstract(abstract)
	if _, ok := container.resolved[index]; ok {
		return true
	}
	if _, ok := container.instances[index]; ok {
		return true
	}
	return false
}

// 解析一个接口
func (container *Container) resolve(abstract interface{}, parameters []interface{}, raiseEvents bool) (object interface{}) {

	index := container.checkAbstract(abstract)

	// todo: 如果是有别名则使用别名
	// todo: 如果是需要上下文条件的时候则 获取上下文条件
	// 如果没有上下文出现了绑定了实例则直接返回
	if _, ok := container.instances[index]; ok {
		return container.instances["concrete"]
	}

	// 如果没有则进行构造
	//concrete := container.getConcrete(abstract)

	//// 开始进行构建
	//if container.isBuildable(abstract, container) {
	//	//object = container.Build(concrete, parameters)
	//}

	// 判定是否为共享的，如果是共享的则进行共享
	if container.IsShared(index) {
		container.instances[index] = object
	}

	if raiseEvents {
		// 判定是否下发事件
	}

	// 记录最后被解析的事件
	container.resolved[index] = true

	return object
}

// 重新绑定一个接口
func (container *Container) rebound(abstract string) interface{} {
	instance := container.Make(abstract)
	// todo 发送一些事件，并且上报被重新 绑定了 ， 一般是不允许重新绑定

	return instance
}

// 获取实际的实现方式
func (container *Container) getConcrete(abstract interface{}) (concrete interface{}) {

	index := container.checkAbstract(abstract)
	// TODO : 如果存在上下文的绑定则返回上下文的内容

	// 如果 设置了绑定的内容则返回绑定的内容
	if value, ok := container.bindings[index]; ok {
		return value.concrete
	}
	concrete = abstract
	return concrete
}

// IsShared 判断这个接口是否为共享的
func (container *Container) IsShared(abstract string) bool {
	if _, ok := container.instances[abstract]; ok {
		return true
	}
	if value, ok := container.bindings[abstract]; ok {
		return value.shared
	}
	return false
}


