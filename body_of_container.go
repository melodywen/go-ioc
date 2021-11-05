package container

type Bind struct {
	shared   bool
	concrete interface{}
}

type BodyOfContainer struct {
	instances map[string]interface{} // 绑定的实例 ， 如果他是单例模式则全部存储到这里面
	bindings  map[string]Bind        // 绑定的策略及其配置
	resolved  map[string]bool        // 是否最终解析成功

	BuildOfContainer
	CommonOfContainer
}

// 用于测试使用
func newBodyOfContainer() *BodyOfContainer {
	instances := map[string]interface{}{}
	bindings := map[string]Bind{}
	resolved := map[string]bool{}
	return &BodyOfContainer{instances: instances, bindings: bindings, resolved: resolved}
}

// DropStaleInstances 移除已经缓存的实例å
func (body *BodyOfContainer) DropStaleInstances(abstract interface{}) bool {
	index := body.AbstractToString(abstract)
	if _, ok := body.instances[index]; ok {
		delete(body.instances, index)
		return true
	}
	return false
}

// Resolved 是否已经实例化过
func (body *BodyOfContainer) Resolved(abstract interface{}) bool {
	index := body.AbstractToString(abstract)
	if _, ok := body.resolved[index]; ok {
		return true
	}
	if _, ok := body.instances[index]; ok {
		return true
	}
	return false
}

// Rebound 再次绑定的操作  todo 待实现
func (body *BodyOfContainer) Rebound(abstract string) {
	//instance := body.Make(abstract)
	// todo 发送一些事件，并且上报被重新 绑定了 ， 一般是不允许重新绑定
	panic("如果重新绑定则需要进行上报")
}

// Bind 绑定接口
func (body *BodyOfContainer) Bind(abstract interface{}, concrete interface{}, shared bool) {
	// 获取对应的 map key
	index := body.AbstractToString(abstract)

	// 删除老旧的实例
	body.DropStaleInstances(index)

	// todo: 如果他不是一个闭包函数，则把它扩展为一个闭包函数，后面扩展就方便很多

	// 直接进行绑定
	body.bindings[index] = Bind{
		shared:   shared,
		concrete: concrete,
	}

	// 如果是之前已经绑定过则再次重新绑定,
	if body.Resolved(abstract) {
		body.Rebound(index)
	}
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
// todo 根据上下文的进行绑定先省略
func (body *BodyOfContainer) GetConcrete(abstract interface{}) (concrete interface{}) {

	index := body.AbstractToString(abstract)
	// TODO : 如果存在上下文的绑定则返回上下文的内容

	// 如果 设置了绑定的内容则返回绑定的内容
	if value, ok := body.bindings[index]; ok {
		return value.concrete
	}
	concrete = abstract
	return concrete
}
