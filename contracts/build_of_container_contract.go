package contracts

type BuildOfContainerContract interface {
	// Build 动态构建一个实例出来：
	Build(concrete interface{}, parameters []interface{}) (object interface{})
}
