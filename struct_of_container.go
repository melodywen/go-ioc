package container

// Bind 绑定的结构体
type Bind struct {
	shared   bool
	concrete interface{}
}

// BuildStackInContainer 当需要依赖输出build 过程中的堆栈信息
type BuildStackInContainer struct {
	buildStack []string
}

// StructOfContainer 容器的结构体
type StructOfContainer struct {
	/**
	主体部分
	*/
	resolved  map[string]bool        // 存储是否解析过
	bindings  map[string]Bind        // 绑定的策略及其配置
	instances map[string]interface{} // 绑定的实例 ， 如果他是单例模式则全部存储到这里面

	/**
	别名
	*/
	aliases         map[string]string   // abstract 对应的别名
	abstractAliases map[string][]string // 存放指定的 abstract别名的集合

	/**
	回调函数
	*/
	extenders        map[string][]func(object interface{}, container *Container) interface{} // 当make 出来的数据做多层装饰器
	reboundCallbacks map[string][]func(container *Container, instance interface{})           // 重新绑定的回调函数

	globalBeforeResolvingCallbacks []func(string, []interface{}, *Container) // 全局的回调函数-前置
	globalResolvingCallbacks       []func(interface{}, *Container)           // 全局的回调函数-调用时候
	globalAfterResolvingCallbacks  []func(interface{}, *Container)           // 全局的回调函数-后置

	beforeResolvingCallbacks map[string][]func(string, []interface{}, *Container) // 具体的接口回调事件——前置
	resolvingCallbacks       map[string][]func(interface{}, *Container)           // 全具体的接口回调事件——调用的时候
	afterResolvingCallbacks  map[string][]func(interface{}, *Container)           // 具体的接口回调事件——后置

	// contextual
	contextual map[string]map[string]interface{} // 存放上下文内容
}
