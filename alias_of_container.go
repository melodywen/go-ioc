package container

import "fmt"

// IsAlias 是否为别名
func (container *Container) IsAlias(abstract interface{}) (ok bool) {
	index := container.AbstractToString(abstract)
	_, ok = container.aliases[index]
	return ok
}

// GetAlias 获取
func (container *Container) GetAlias(abstract interface{}) string {
	index := container.AbstractToString(abstract)
	if _, ok := container.aliases[index]; ok {
		return container.GetAlias(container.aliases[index])
	}
	return index
}

// Alias 设置别名
func (container *Container) Alias(abstract interface{}, alias interface{}) {
	abstractStr := container.AbstractToString(abstract)
	aliasStr := container.AbstractToString(alias)
	if aliasStr == abstractStr {
		// todo 抛错异常得修改
		panic("[{$abstract}] is aliased to itself.")
	}
	container.aliases[aliasStr] = abstractStr
	if container.abstractAliases[abstractStr] == nil {
		container.abstractAliases[abstractStr] = []string{}
	}
	container.abstractAliases[abstractStr] = append(container.abstractAliases[abstractStr], aliasStr)
}

//  Remove an alias from the contextual binding alias cache.
func (container *Container) removeAbstractAlias(search string) (ok bool) {
	if _, ok := container.aliases[search]; !ok {
		return false
	}
	for abstract, aliases := range container.abstractAliases {
		for index, alias := range aliases {
			fmt.Println(alias, index)
			if alias == search {
				container.abstractAliases[abstract] = append(aliases[:index], aliases[index+1:]...)
				ok = true
			}
		}
	}
	return ok
}
