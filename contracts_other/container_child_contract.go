package contracts

// ContainerChildContract
// @Description: 容器子类必须集成
type ContainerChildContract interface {
	ResolveCallback(identifier string)
}
