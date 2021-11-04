package contracts

type BuildOfContainerContract interface{
	Build(concrete interface{}, parameters []interface{}) (object interface{}) // 进行构建

}
