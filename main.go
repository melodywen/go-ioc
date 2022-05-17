package main

import (
	"fmt"
	"reflect"
)

type AT interface {
	Say() string
}

type animal struct {
	name string
}

func (an *animal) Say() string {
	return an.name
}

func main() {
	var aa AT

	aa = &animal{name: "sdf"}

	fmt.Println(reflect.TypeOf(aa).Kind())
	fmt.Println(reflect.ValueOf(aa).Kind())
	fmt.Println(123131)

	concreteType := reflect.ValueOf(abc).Type()
	fmt.Println(concreteType)
	param := concreteType.In(0)

	fmt.Println(param, 3333)
	fmt.Println(param.Kind(), 3333)
	fmt.Println(param.Name(), 3333)
	fmt.Println("--------------------")
	var bb *AT

	fmt.Println(bb, 3333)
	fmt.Println(reflect.TypeOf(bb).Elem() == param, 3333)

}

func abc(at AT) {

}
