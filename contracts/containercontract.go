package contracts

import (
	"fmt"
	"reflect"
)

type ContainerContracts interface {
	Bind(abstract interface{}, concrete interface{},shared bool)
	DropStaleInstances(abstract interface{})

}

func aa()  {
	fmt.Println(reflect.Interface)

}