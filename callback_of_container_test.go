package container

import (
	"cjw.com/melodywen/go-ioc/mock"
	"reflect"
	"testing"
)

func TestContainer_Rebinding(t *testing.T) {
	type fields struct {
		flagName string
		flagNum  int
	}
	type args struct {
		abstract interface{}
		callback func(container *Container, instance interface{})
	}
	tests := []struct {
		fields fields
		args   args
		want   interface{}
	}{
		{
			fields: fields{
				flagName: "如果是第一次绑定",
				flagNum:  0,
			},
		},
		{
			fields: fields{
				flagName: "如果再次绑定未发现内容",
				flagNum:  1,
			},
		},
	}
	container := newContainer()
	for _, tt := range tests {
		t.Run(tt.fields.flagName, func(t *testing.T) {
			switch tt.fields.flagNum {
			case 0:
				container.Rebinding(mock.Animal{}, func(container *Container, instance interface{}) {
					if _, ok := container.reboundCallbacks[container.AbstractToString(mock.Animal{})]; !ok {
						t.Errorf("rebinding 有问题")
					}
				})
				container.Bind(mock.Animal{}, func() *mock.Animal {
					return mock.NewAnimal("dog", 12, "cate")
				}, true)
				if !reflect.DeepEqual(container.Make(mock.Animal{}), mock.NewAnimal("dog", 12, "cate")) {
					t.Errorf("make 的内容有误")
				}
				container.Rebinding(mock.Animal{}, func(container *Container, instance interface{}) {
					if _, ok := container.reboundCallbacks[container.AbstractToString(mock.Animal{})]; !ok {
						t.Errorf("rebinding 有问题")
					}
				})
				container.Bind(mock.Animal{}, func() *mock.Animal {
					return mock.NewAnimal("dog", 12, "cate")
				}, true)
				if !reflect.DeepEqual(container.Make(mock.Animal{}), mock.NewAnimal("dog", 12, "cate")) {
					t.Errorf("make 的内容有误")
				}
			case 1:
				container.Bind(mock.Animal{}, func() *mock.Animal {
					return mock.NewAnimal("dog", 12, "cate")
				}, true)
				container.Bind(mock.Animal{}, func() *mock.Animal {
					return mock.NewAnimal("dog", 12, "cate")
				}, true)
			}

		})
	}
}

func TestContainer_Extend(t *testing.T) {
	type fields struct {
		flagName string
		flagNum  int
	}
	type args struct {
		abstract string
		closure  func(object interface{}, container *Container) interface{}
	}
	tests := []struct {
		fields fields
		args   args
	}{
		{
			fields: fields{
				flagName: "测试主体流程，并且一个extend",
				flagNum:  0,
			},
		},
		{
			fields: fields{
				flagName: "测试主体流程，并且多个extend",
				flagNum:  1,
			},
		},
		{
			fields: fields{
				flagName: "测试主体流程，forgetExtenders",
				flagNum:  2,
			},
		},
		{
			fields: fields{
				flagName: "测试主体流程，绑定instance",
				flagNum:  3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.fields.flagName, func(t *testing.T) {
			container := newContainer()

			switch tt.fields.flagNum {
			case 0:
				container.Bind(mock.Animal{}, func() (int, int, int) {
					return mock.AddAndParam(3, 5)
				}, false)
				result := container.MakeWithParams(mock.Animal{}, []interface{}{})
				container.Extend(mock.Animal{}, func(object interface{}, container *Container) interface{} {
					if !reflect.DeepEqual(object, []interface{}{8, 3, 5}) {
						t.Errorf("第一层有问题")
					}
					var response []interface{}
					for _, v := range object.([]interface{}) {
						response = append(response, v.(int)+3)
					}
					return response
				})
				result = container.MakeWithParams(mock.Animal{}, []interface{}{})
				if !reflect.DeepEqual(result, []interface{}{11, 6, 8}) {
					t.Errorf("第一层有问题555")
				}

			case 1:
				for i := 1; i < 10; i++ {
					container.Extend(mock.Animal{}, func(object interface{}, container *Container) interface{} {
						var response []interface{}
						for _, v := range object.([]interface{}) {
							response = append(response, v.(int)+3)
						}
						return response
					})
				}
				container.Bind(mock.Animal{}, mock.AddAndParam, false)
				result := container.MakeWithParams(mock.Animal{}, []interface{}{3, 5})

				if !reflect.DeepEqual(result, []interface{}{35, 30, 32}) {
					t.Errorf("第10 层有问题")
				}
			case 2:
				container.Extend(mock.Animal{}, func(object interface{}, container *Container) interface{} {
					var response []interface{}
					for _, v := range object.([]interface{}) {
						response = append(response, v.(int)+3)
					}
					return response
				})
				container.ForgetExtenders(&mock.Animal{})
				container.ForgetExtenders(mock.Animal{})
			case 3:
				container.Instance(mock.Animal{}, []interface{}{1, 2, 3})
				container.Extend(mock.Animal{}, func(object interface{}, container *Container) interface{} {
					var response []interface{}
					for _, v := range object.([]interface{}) {
						response = append(response, v.(int)+3)
					}
					return response
				})
				result := container.Make(mock.Animal{})
				if !reflect.DeepEqual(result, []interface{}{4, 5, 6}) {
					t.Errorf("解析实例有误")
				}
			}
		})
	}
}
