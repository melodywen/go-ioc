package container

import (
	"cjw.com/melodywen/go-ioc/mock"
	"fmt"
	"reflect"
	"testing"
)

func TestContainer_dropStaleInstances(t *testing.T) {
	type fields struct {
		StructOfContainer
	}
	type args struct {
		abstract string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		{
			name: "测试移除实例->true",
			fields: fields{
				StructOfContainer{
					instances: map[string]interface{}{"aa": 1},
					aliases:   map[string]string{"ab": "aa"},
				},
			},
			args:   args{abstract: "ab"},
			wantOk: true,
		}, {
			name: "测试移除实例->false",
			fields: fields{
				StructOfContainer{
					instances: map[string]interface{}{"aa": 1},
					aliases:   map[string]string{"ab": "aa"},
				},
			},
			args:   args{abstract: "bb"},
			wantOk: false,
		}, {
			name: "测试移除实例->全部移除",
			fields: fields{
				StructOfContainer{
					instances: map[string]interface{}{"aa": 1},
					aliases:   map[string]string{"aa": "aa"},
				},
			},
			args:   args{abstract: "aa"},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := &Container{
				StructOfContainer: tt.fields.StructOfContainer,
			}
			if gotOk := container.dropStaleInstances(tt.args.abstract); gotOk != tt.wantOk {
				t.Errorf("dropStaleInstances() = %v, want %v", gotOk, tt.wantOk)
			} else {
				fmt.Println(container)
			}
		})
	}
}

func TestContainer_Resolved(t *testing.T) {
	type fields struct {
		StructOfContainer
	}
	type args struct {
		abstract interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		{
			name: "如果在instance 中",
			fields: fields{StructOfContainer{
				instances: map[string]interface{}{"cjw.com/melodywen/go-ioc/mock.Animal": 1},
			}},
			args:   args{abstract: mock.Animal{}},
			wantOk: true,
		},
		{
			name: "如果在resolved 中",
			fields: fields{StructOfContainer{
				resolved: map[string]bool{"*cjw.com/melodywen/go-ioc/mock.Animal": true},
			}},
			args:   args{abstract: &mock.Animal{}},
			wantOk: true,
		},
		{
			name: "如果都不在中",
			fields: fields{StructOfContainer{
				instances: map[string]interface{}{"cjw.com/melodywen/go-ioc/mock.Animal": 1},
				resolved:  map[string]bool{"*cjw.com/melodywen/go-ioc/mock.Animal": true},
			}},
			args:   args{abstract: 1},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := &Container{
				StructOfContainer: tt.fields.StructOfContainer,
			}
			if gotOk := container.Resolved(tt.args.abstract); gotOk != tt.wantOk {
				t.Errorf("Resolved() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestContainer_Bound(t *testing.T) {
	type fields struct {
	}
	type args struct {
		abstract interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		{
			name: "测试",
		},
	}
	container := newContainer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container.Bind(mock.Animal{}, mock.NewAnimal, false)
			if gotOk := container.Bound(mock.Animal{}); gotOk != true {
				t.Errorf("Bound() = %v, want %v", gotOk, true)
			}
			container.Instance(mock.NewAnimalAndParam, func() *mock.Animal {
				return mock.NewAnimal("dog", 12, "cate")
			})
			//container.MakeWithParams(mock.NewAnimalAndParam, []interface{}{})
			if gotOk := container.Bound(mock.NewAnimalAndParam); gotOk != true {
				t.Errorf("Bound() = %v, want %v", gotOk, tt.wantOk)
			}
			container.Alias(mock.Animal{}, "aabb")
			if gotOk := container.Bound("aabb"); gotOk != true {
				t.Errorf("Bound() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestContainer_Instance(t *testing.T) {
	type fields struct {
		abstract interface{}
		instance interface{}
	}
	type args struct {
		abstract interface{}
		instance interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		// TODO: Add test cases.
		{
			name: "测试如果在别名 和 bind 中则都进行删除",
			want: mock.NewAnimal("dog", 12, "cat"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := newContainer()
			container.Bind(mock.Animal{}, func() *mock.Animal {
				return mock.NewAnimal("dog", 12, "cat")
			}, true)
			container.Alias(mock.Animal{}, mock.NewAnimalAndParam)
			container.MakeWithParams(mock.Animal{}, []interface{}{})
			//fmt.Println(container)
			got := container.Instance(mock.NewAnimalAndParam, mock.NewAnimal("dog", 12, "cat"))
			//fmt.Println(container)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Instance() = %v, want %v", got, tt.want)
			}
		})
	}
}
