package container

import (
	contracts2 "github.com/melodywen/go-ioc/contracts_other"
	"sync"
)

// bind
//  @Description: 绑定的原始信息
type bind struct {
	shared   bool
	concrete any
}

// containerStack
//  @Description: 每次进行解析的时候进行压堆栈上下文信息
type containerStack struct {
	Stack []string
}

// newContainerStack
//  @Description:
//  @return *containerStack
func newContainerStack() *containerStack {
	return &containerStack{}
}

// containerStore
//  @Description:存储绑定相应信息并且记录实例化信息
type containerStore struct {
	storeMutex sync.Mutex      // 互斥锁
	resolved   map[string]bool // 存储是否解析过
	bindings   map[string]bind // 绑定的策略及其配置
	instances  map[string]any  // 绑定的实例 ， 如果他是单例模式则全部存储到这里面
}

// containerAliases
//  @Description:别名
type containerAliases struct {
	aliasesMutex    sync.Mutex          // 互斥锁
	aliases         map[string]string   // abstract 对应的别名
	abstractAliases map[string][]string // 存放指定的 abstract别名的集合
}

// containerCallback
//  @Description: 回调函数
type containerCallback struct {
	callbackMutex    sync.Mutex                                              // 互斥锁
	extenders        map[string][]func(object any, container *Container) any // 当make 出来的数据做多层装饰器
	reboundCallbacks map[string][]func(container *Container, instance any)   // 重新绑定的回调函数

	globalBeforeResolvingCallbacks []func(string, []any, *Container) // 全局的回调函数-前置
	globalResolvingCallbacks       []func(any, *Container)           // 全局的回调函数-调用时候
	globalAfterResolvingCallbacks  []func(any, *Container)           // 全局的回调函数-后置

	beforeResolvingCallbacks map[string][]func(string, []any, *Container) // 具体的接口回调事件——前置
	resolvingCallbacks       map[string][]func(any, *Container)           // 全具体的接口回调事件——调用的时候
	afterResolvingCallbacks  map[string][]func(any, *Container)           // 具体的接口回调事件——后置
}

// containerContext
//  @Description: 上下文
type containerContext struct {
	contextual map[string]map[string]any // 存放上下文内容
}

// Container
//  @Description: 容器实例
type Container struct {
	containerStore
	containerAliases
	containerCallback
	containerContext
	child contracts2.ContainerChildContract
}

// NewContainer
//  @Description: 实例化
//  @return *Container
func NewContainer() *Container {
	return &Container{
		containerStore: containerStore{
			storeMutex: sync.Mutex{},
			resolved:   map[string]bool{},
			bindings:   map[string]bind{},
			instances:  map[string]any{},
		},
		containerAliases: containerAliases{
			aliasesMutex:    sync.Mutex{},
			aliases:         map[string]string{},
			abstractAliases: map[string][]string{},
		},
		containerCallback: containerCallback{
			extenders:        map[string][]func(object any, container *Container) any{},
			reboundCallbacks: map[string][]func(container *Container, instance any){},

			globalBeforeResolvingCallbacks: []func(string, []any, *Container){},
			globalResolvingCallbacks:       []func(any, *Container){},
			globalAfterResolvingCallbacks:  []func(any, *Container){},

			beforeResolvingCallbacks: map[string][]func(string, []any, *Container){},
			resolvingCallbacks:       map[string][]func(any, *Container){},
			afterResolvingCallbacks:  map[string][]func(any, *Container){},
		},
		containerContext: containerContext{contextual: map[string]map[string]any{}},
		child:            nil,
	}
}
