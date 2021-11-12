package container

// IsAlias 是否为别名
func (container *Container) IsAlias(abstract interface{}) (ok bool) {
	index := container.AbstractToString(abstract)
	_, ok = container.aliases[index]
	return ok
}

// GetAlias 获取
func (container *Container)GetAlias(abstract interface{}) string {
	index := container.AbstractToString(abstract)
	if _, ok := container.aliases[index]; ok {
		return container.GetAlias(container.aliases[index])
	}
	return index
}
