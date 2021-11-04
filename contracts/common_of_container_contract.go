package contracts

type CommonOfContainerContract interface {
	// AbstractToString 通过 abstract 进行 字符串得到 作为map key
	AbstractToString(abstract interface{}) string
}
