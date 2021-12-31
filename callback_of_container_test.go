package container

import (
	"fmt"
	"github.com/melodywen/go-ioc/mock"
	"reflect"
	"strconv"
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
	container := NewContainer()
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
			container := NewContainer()

			switch tt.fields.flagNum {
			case 0:
				container.Bind(mock.Animal{}, func() (int, int, int) {
					return mock.AddAndParam(3, 5)
				}, false)
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
				result := container.MakeWithParams(mock.Animal{}, []interface{}{})
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

func TestContainer_BeforeResolving(t *testing.T) {
	type fields struct {
		flagName string
		flagNum  int
		want     []string
	}

	tests := []struct {
		fields fields
	}{
		{
			fields: fields{
				flagName: "第一个简单的回调函数",
				flagNum:  0,
				want: []string{
					"全局before回调函数",
					"全局before回调函数2",
					"私下的回调before回调函数1",
					"私下的回调before回调函数2",
				},
			},
		},
	}
	container := NewContainer()
	for _, tt := range tests {
		t.Run(tt.fields.flagName, func(t *testing.T) {
			switch tt.fields.flagNum {
			case 0:
				got := []string{}
				container.Bind(mock.Animal{}, mock.NewAnimal, true)
				container.BeforeResolving(nil, func(s string, i []interface{}, container *Container) {
					got = append(got, "全局before回调函数")
					fmt.Println(tt.fields.flagNum, got, s)
				})
				container.BeforeResolving(nil, func(s string, i []interface{}, container *Container) {
					got = append(got, "全局before回调函数2")
					fmt.Println(tt.fields.flagNum, got, s)
				})
				container.BeforeResolving(mock.Animal{}, func(s string, i []interface{}, container *Container) {
					got = append(got, "私下的回调before回调函数1")
					fmt.Println(tt.fields.flagNum, got, s)
				})
				container.BeforeResolving(mock.Animal{}, func(s string, i []interface{}, container *Container) {
					got = append(got, "私下的回调before回调函数2")
					fmt.Println(tt.fields.flagNum, got, s)
				})
				container.MakeWithParams(mock.Animal{}, []interface{}{"tom", 3, "cate"})
				if !reflect.DeepEqual(tt.fields.want, got) {
					t.Errorf("解析实例有误")
				}
			}
		})
	}
}

func TestContainer_Resolving(t *testing.T) {
	type fields struct {
		flagName string
		flagNum  int
		want     []string
	}

	tests := []struct {
		fields fields
	}{
		{
			fields: fields{
				flagName: "绑定一个-resolve",
				flagNum:  0,
				want: []string{
					"0-global-resolve-callback-01", "0-global-resolve-callback-02",
					"0-local-resolve-callback-01", "0-local-resolve-callback-02",
				},
			},
		}, {
			fields: fields{
				flagName: "绑定其他的-resolve",
				flagNum:  1,
				want: []string{
					"1-global-resolve-callback-01", "1-global-resolve-callback-02", "1-global-resolve-callback-11",
					"1-global-resolve-callback-12", "1-local-resolve-callback-11", "1-local-resolve-callback-11",
				},
			},
		},
		{
			fields: fields{
				flagName: "绑定一个-after-resolve",
				flagNum:  2,
				want: []string{
					"2-global-resolve-callback-01", "2-global-resolve-callback-02", "2-global-resolve-callback-11",
					"2-global-resolve-callback-12", "2-local-resolve-callback-01", "2-local-resolve-callback-02",
					"2-global-after-resolve-callback-21", "2-global-after-resolve-callback-22",
					"2-local-after-resolve-callback-21", "2-local-after-resolve-callback-22",
				},
			},
		},
	}
	container := NewContainer()
	got := []string{}
	for _, tt := range tests {
		t.Run(tt.fields.flagName, func(t *testing.T) {
			switch tt.fields.flagNum {
			case 0:
				container.Bind(mock.Animal{}, func() *mock.Animal {
					return mock.NewAnimal("cat", 12, "cate")
				}, true)
				container.Resolving(nil, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-global-resolve-callback-01")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.Resolving(nil, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-global-resolve-callback-02")
					fmt.Println(tt.fields.flagNum, got, object)
				})

				container.Resolving(mock.Animal{}, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-local-resolve-callback-01")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.Resolving(mock.Animal{}, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-local-resolve-callback-02")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.MakeWithParams(mock.Animal{}, []interface{}{})
				container.MakeWithParams(mock.Animal{}, []interface{}{})
				//container.MakeWithParams(mock.Animal{},[]interface{}{})
				if !reflect.DeepEqual(tt.fields.want, got) {
					fmt.Println(tt.fields.want, got)
					t.Errorf("解析实例有误")
				}
			case 1:
				container.Bind(mock.NewAnimalAndParam, func() *mock.Animal {
					return mock.NewAnimal("cat", 12, "cate")
				}, false)
				container.Resolving(nil, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-global-resolve-callback-11")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.Resolving(nil, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-global-resolve-callback-12")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.Resolving(mock.NewAnimalAndParam, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-local-resolve-callback-11")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.Resolving(mock.NewAnimalAndParam, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-local-resolve-callback-11")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.MakeWithParams(mock.NewAnimalAndParam, []interface{}{})
				if !reflect.DeepEqual(tt.fields.want, got) {
					fmt.Println(tt.fields.want, got)
					t.Errorf("解析实例有误")
				}
			case 2:
				container.Bind(mock.Animal{}, func() *mock.Animal {
					return mock.NewAnimal("cat", 12, "cate")
				}, false)
				got = []string{}
				container.AfterResolving(nil, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-global-after-resolve-callback-21")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.AfterResolving(nil, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-global-after-resolve-callback-22")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.AfterResolving(mock.Animal{}, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-local-after-resolve-callback-21")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.AfterResolving(mock.Animal{}, func(object interface{}, container *Container) {
					got = append(got, strconv.Itoa(tt.fields.flagNum)+"-local-after-resolve-callback-22")
					fmt.Println(tt.fields.flagNum, got, object)
				})
				container.MakeWithParams(mock.Animal{}, []interface{}{})
				if !reflect.DeepEqual(tt.fields.want, got) {
					fmt.Println(tt.fields.want, got)
					t.Errorf("解析实例有误")
				}
			}
			got = []string{}
		})
	}
}
