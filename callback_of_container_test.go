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
	}
	container := newContainer()
	for _, tt := range tests {
		t.Run(tt.fields.flagName, func(t *testing.T) {
			switch tt.fields.flagNum {
			case 0:
				container.Bind(mock.Animal{}, func() *mock.Animal {
					return mock.NewAnimal("dog", 12, "cate")
				}, true)
				if !reflect.DeepEqual(container.Make(mock.Animal{}), mock.NewAnimal("dog", 12, "cate")) {
					t.Errorf("make 的内容有误")
				}
				container.Rebinding(mock.Animal{}, func(container *Container, instance interface{}) {
					if _,ok := container.reboundCallbacks[container.AbstractToString(mock.Animal{})]; !ok {
						t.Errorf("rebinding 有问题")
					}
				})
				container.Bind(mock.Animal{}, func() *mock.Animal {
					return mock.NewAnimal("dog", 12, "cate")
				}, true)
				if !reflect.DeepEqual(container.Make(mock.Animal{}), mock.NewAnimal("dog", 12, "cate")) {
					t.Errorf("make 的内容有误")
				}

			}

		})
	}
}
