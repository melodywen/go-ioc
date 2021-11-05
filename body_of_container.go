package container

import (
	"fmt"
	"reflect"
)

type Bind struct {
	shared   bool
	concrete interface{}
}

type BodyOfContainer struct {
	instances map[string]interface{} // 绑定的实例 ， 如果他是单例模式则全部存储到这里面
	bindings  map[string]Bind        // 绑定的策略及其配置
	resolved  map[string]bool        // 是否最终解析成功

	aliases          map[string]string        // abstract 对应的别名
	reboundCallbacks map[string][]interface{} // 重新绑定的回调函数
	extenders        map[string][]interface{} // 当make 出来的数据做多层装饰器

	BuildOfContainer
	CommonOfContainer
}

// 用于测试使用
func newBodyOfContainer() *BodyOfContainer {
	return &BodyOfContainer{
		instances: map[string]interface{}{},
		bindings:  map[string]Bind{},
		resolved:  map[string]bool{},
		aliases:   map[string]string{},
	}
}

// DropStaleInstances 移除已经缓存的实例 和别名
func (body *BodyOfContainer) DropStaleInstances(abstract string) {
	if _, ok := body.instances[abstract]; ok {
		delete(body.instances, abstract)
	}
	if _, ok := body.aliases[abstract]; ok {
		delete(body.aliases, abstract)
	}
}

// IsAlias 是否为别名
func (body *BodyOfContainer) IsAlias(abstract interface{}) bool {
	index := body.AbstractToString(abstract)
	if _, ok := body.aliases[index]; ok {
		return true
	}
	return false
}

// GetAlias 获取
func (body *BodyOfContainer) GetAlias(abstract interface{}) string {
	index := body.AbstractToString(abstract)
	if _, ok := body.aliases[index]; ok {
		return body.GetAlias(body.aliases[index])
	}
	return index
}

// Resolved 是否已经实例化过
func (body *BodyOfContainer) Resolved(abstract interface{}) bool {
	if body.IsAlias(abstract) {
		abstract = body.GetAlias(abstract)
	}
	index := body.AbstractToString(abstract)

	if _, ok := body.resolved[index]; ok {
		return true
	}
	if _, ok := body.instances[index]; ok {
		return true
	}
	return false
}

// Rebound 再次绑定的操作
func (body *BodyOfContainer) Rebound(abstract string) {
	instance := body.Make(abstract)
	for callback := range body.getReboundCallbacks(abstract) {
		if reflect.TypeOf(callback).Kind() != reflect.Func {
			continue
		}
		fmt.Println(instance)
		panic("开始绑定回调")
	}
}

func (body *BodyOfContainer) getReboundCallbacks(abstract string) (reboundCallbacks []interface{}) {
	ok := false
	if reboundCallbacks, ok = body.reboundCallbacks[abstract]; ok {
		return reboundCallbacks
	}
	return []interface{}{}
}

func (body *BodyOfContainer) getExtenders(abstract string) (extenderCallbacks []interface{}) {
	ok := false
	if extenderCallbacks, ok = body.extenders[abstract]; ok {
		return extenderCallbacks
	}
	return []interface{}{}
}

// Bind 绑定接口
func (body *BodyOfContainer) Bind(abstract interface{}, concrete interface{}, shared bool) {
	// 获取对应的 map key
	index := body.AbstractToString(abstract)

	// 删除老旧的实例
	body.DropStaleInstances(index)

	if concrete == nil {
		concrete = abstract
	}

	if reflect.TypeOf(concrete).Kind() != reflect.Func {
		concrete = body.getClosure(index, concrete)
	}

	// 直接进行绑定
	body.bindings[index] = Bind{shared: shared, concrete: concrete}

	// 如果是之前已经绑定过则再次重新绑定,
	if body.Resolved(abstract) {
		body.Rebound(index)
	}
}

func (body *BodyOfContainer) getClosure(abstract string, concrete interface{}) func(body *BodyOfContainer) interface{} {
	return func(body *BodyOfContainer) interface{} {
		if reflect.DeepEqual(abstract, concrete) {
			return body.Build(concrete, []interface{}{})
		}
		return body.resolve(abstract, []interface{}{}, false)
	}
}

// Make 对外暴露make 方法
func (body *BodyOfContainer) Make(abstract interface{}) interface{} {

	return 123
}

// MakeWithParams 对外暴露make 方法并且
func (body *BodyOfContainer) MakeWithParams(abstract interface{}, parameters []interface{}) interface{} {
	index := body.AbstractToString(abstract)
	return body.resolve(index, parameters, true)
}

// IsShared 判断这个接口是否为共享
func (body *BodyOfContainer) IsShared(abstract string) bool {
	if _, ok := body.instances[abstract]; ok {
		return true
	}
	if value, ok := body.bindings[abstract]; ok {
		return value.shared
	}
	return false
}

// GetConcrete 通过接口获取具体实现
func (body *BodyOfContainer) GetConcrete(abstract string) (concrete interface{}) {
	// 如果上下环境中有内容
	if concrete := body.getContextualConcrete(abstract); concrete != nil {
		return concrete
	}

	// 如果 设置了绑定的内容则返回绑定的内容
	if value, ok := body.bindings[abstract]; ok {
		return value.concrete
	}
	return abstract
}
// TODO 这一部分 先不实现，因为设计到并发 锁的问题
func (body *BodyOfContainer) getContextualConcrete(abstract string) (concrete interface{}) {

	return nil
}

// todo 这一部分 先不实现， 进行回调上报
func (body *BodyOfContainer) fireResolvingCallbacks(abstract string, object interface{}) {

}

// 解析一个接口
func (body *BodyOfContainer) resolve(abstract string, parameters []interface{}, raiseEvents bool) (object interface{}) {
	abstract = body.GetAlias(abstract)

	needsContextualBuild := len(parameters) != 0 || body.getContextualConcrete(abstract) != nil

	if _, ok := body.instances[abstract]; ok && !needsContextualBuild {
		return body.instances[abstract]
	}

	concrete := body.GetConcrete(abstract)

	// 如果可以构建则直接构建 否则递归判定
	if body.IsBuildable(abstract, concrete) {
		object = body.Build(concrete, parameters)
	} else {
		object = body.MakeWithParams(concrete, parameters)
	}

	// 做装饰器回调
	for callback := range body.getExtenders(abstract) {
		if reflect.TypeOf(callback).Kind() != reflect.Func {
			continue
		}
		fmt.Println(callback)
		panic("开始绑定回调")
	}

	if body.IsShared(abstract) {
		body.instances[abstract] = object
	}

	if raiseEvents {
		// 判定是否下发事件
		body.fireResolvingCallbacks(abstract, object)
	}

	// 记录最后被解析的事件
	body.resolved[abstract] = true

	return object
}
