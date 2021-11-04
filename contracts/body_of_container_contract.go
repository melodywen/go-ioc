package contracts

type BodyOfContainerContract interface {
	// DropStaleInstances 移除已经缓存的实例å
	DropStaleInstances(abstract interface{}) bool
	// Resolved 是否已经实例化过
	Resolved(abstract interface{}) bool
	// Rebound 再次绑定的操作
	Rebound(abstract string) bool

	// IsShared 判断这个接口是否为共享
	IsShared(abstract string) bool
	// GetConcrete 通过接口获取具体实现
	GetConcrete(abstract interface{}) (concrete interface{})

	CommonOfContainerContract
	BuildOfContainerContract
}
