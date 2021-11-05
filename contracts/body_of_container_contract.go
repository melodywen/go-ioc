package contracts

type BodyOfContainerContract interface {
	// DropStaleInstances 移除已经缓存的实例和别名
	DropStaleInstances(abstract string)

	// IsAlias 是否为别名
	IsAlias(abstract interface{}) bool
	// GetAlias 获取
	GetAlias(abstract interface{}) interface{}
	// Resolved 是否已经实例化过
	Resolved(abstract interface{}) bool
	// Rebound 再次绑定的操作
	Rebound(abstract string)
	// Bind 绑定接口
	Bind(abstract interface{}, concrete interface{}, shared bool)

	// IsShared 判断这个接口是否为共享
	IsShared(abstract string) bool
	// GetConcrete 通过接口获取具体实现
	GetConcrete(abstract interface{}) (concrete interface{})

	CommonOfContainerContract
	BuildOfContainerContract
}
