package container

import (
	"fmt"
	"reflect"
)

type BuildOfContainer struct {
}

// Build 动态构建一个实例出来：
// 1. 可以是某个类的构造方法
// 2. 也可以是回调函数
// 3. 一个具体的标量值
// 4. 如果是一个结构体，则直接寻找他对应的实例化方法
// todo:
// golang 目前没有发现动态加载功能， 待实现的一个功能： 如果是是  concrete 是一个结构体，能否自动寻路找到他的实例化方法？
func (_ *BuildOfContainer) Build(concrete interface{}, parameters []interface{}) (object interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("build 发现异常")
			fmt.Println("concrete:", reflect.TypeOf(concrete))
			fmt.Println("parameters:", parameters)
			panic(err)
		}
	}()
	// 获取实现类的类型
	concreteType := reflect.TypeOf(concrete)

	switch concreteType.Kind() {
	case reflect.Func: // 如果是一个回调函数或者是这个构造方法都走这里
		// 获取实现类的值
		concreteValue := reflect.ValueOf(concrete)
		// 函数的形参绑定
		var params []reflect.Value
		for _, parameter := range parameters {
			params = append(params, reflect.ValueOf(parameter))
		}
		// 调用函数
		resultList := concreteValue.Call(params)
		// 然后进行克隆反射
		numOut := concreteValue.Type().NumOut()
		response := []interface{}{}
		for m := 0; m < numOut; m++ {
			returnType := concreteValue.Type().Out(m)
			switch returnType.Kind() {
			case reflect.Ptr: //如果是指针类型
				returnNew := reflect.New(returnType.Elem()).Elem() //创建对象 //获取源实际类型(否则为指针类型)
				returnValue := resultList[m]                       //源数据值
				returnValue = returnValue.Elem()                   //源数据实际值（否则为指针）
				returnNew.Set(returnValue)                         //设置数据
				returnNew = returnNew.Addr()                       //创建对象的地址（否则返回值）
				response = append(response, returnNew.Interface()) //返回地址
			default:
				returnNew := reflect.New(returnType).Elem()        //创建对象
				returnValue := resultList[m]                       //源数据值
				returnNew.Set(returnValue)                         //设置数据
				response = append(response, returnNew.Interface()) //返回
			}
		}
		if len(response) == 1 {
			return response[0]
		}
		return response
	default:
		fmt.Println("this concrete is :", concrete)
		panic("Target [xxx] is not instantiable. please see upper logger")
	}
}

// isBuildable 判断可以进行构建
// 1. 如果是自己等于自己则直接输出
// 2. 如果是回调函数则保留
// 3. 如果是其他的类型，说明还不够，需要进行递归获取
func (_ *BuildOfContainer) isBuildable(abstract interface{}, concrete interface{}) bool {

	if reflect.DeepEqual(abstract, concrete) {
		return true
	}
	return reflect.TypeOf(concrete).Kind() == reflect.Func
}
