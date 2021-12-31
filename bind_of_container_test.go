package container

import (
	"fmt"
	"github.com/melodywen/go-ioc/mock"
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
		alias    interface{}
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
				instances: map[string]interface{}{"github.com/melodywen/go-ioc/mock.Animal": 1},
			}},
			args:   args{abstract: mock.Animal{}},
			wantOk: true,
		},
		{
			name: "如果在resolved 中",
			fields: fields{StructOfContainer{
				resolved: map[string]bool{"*github.com/melodywen/go-ioc/mock.Animal": true},
			}},
			args:   args{abstract: &mock.Animal{}},
			wantOk: true,
		},
		{
			name: "如果都不在中",
			fields: fields{StructOfContainer{
				instances: map[string]interface{}{"github.com/melodywen/go-ioc/mock.Animal": 1},
				resolved:  map[string]bool{"*github.com/melodywen/go-ioc/mock.Animal": true},
			}},
			args:   args{abstract: 1},
			wantOk: false,
		}, {
			name: "如果是别名",
			fields: fields{StructOfContainer{
				instances:       map[string]interface{}{},
				resolved:        map[string]bool{"abc": true},
				aliases:         map[string]string{},
				abstractAliases: map[string][]string{},
			}},
			args: args{
				abstract: "abc",
				alias:    mock.Animal{},
			},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := &Container{
				StructOfContainer: tt.fields.StructOfContainer,
			}
			var gotOk bool
			if tt.args.alias != nil {
				fmt.Println(tt.args.abstract, tt.args.alias)
				container.Alias(tt.args.abstract, tt.args.alias)
				gotOk = container.Resolved(tt.args.alias)
			} else {
				gotOk = container.Resolved(tt.args.abstract)
			}
			if gotOk != tt.wantOk {
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
	container := NewContainer()
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
		{
			name: "测试如果在别名 和 bind 中则都进行删除",
			want: mock.NewAnimal("dog", 12, "cat"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewContainer()
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

func TestContainer_BindIf(t *testing.T) {

	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
		params   []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "第一次绑定",
			args: args{
				abstract: mock.Animal{},
				concrete: mock.NewAnimal,
				shared:   true,
				params:   []interface{}{"dog", 12, "cate-dog"},
			},
			want: mock.NewAnimal("dog", 12, "cate-dog"),
		}, {
			name: "第二次绑定",
			args: args{
				abstract: mock.Animal{},
				concrete: mock.NewChild,
				shared:   false,
				params:   []interface{}{"dog", 12, "cate-dog"},
			},
			want: mock.NewAnimal("dog", 12, "cate-dog"),
		},
	}
	container := NewContainer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container.BindIf(tt.args.abstract, tt.args.concrete, tt.args.shared)
			got := container.MakeWithParams(tt.args.abstract, tt.args.params)
			fmt.Println(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Instance() = %v, want %v", container.StructOfContainer, tt.want)
			}
		})
	}
}

func TestContainer_SingletonIf(t *testing.T) {
	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
		params   []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "第一次绑定",
			args: args{
				abstract: mock.Animal{},
				concrete: mock.NewAnimal,
				shared:   true,
				params:   []interface{}{"dog", 12, "cate-dog"},
			},
			want: mock.NewAnimal("dog", 12, "cate-dog"),
		}, {
			name: "第二次绑定",
			args: args{
				abstract: mock.Animal{},
				concrete: mock.NewChild,
				shared:   false,
				params:   []interface{}{"dog", 12, "cate-dog"},
			},
			want: mock.NewAnimal("dog", 12, "cate-dog"),
		},
	}
	container := NewContainer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container.SingletonIf(tt.args.abstract, tt.args.concrete)
			got := container.MakeWithParams(tt.args.abstract, tt.args.params)
			fmt.Println(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Instance() = %v, want %v", container.StructOfContainer, tt.want)
			}
		})
	}
}

func TestContainer_Flush(t *testing.T) {
	type fields struct {
		StructOfContainer StructOfContainer
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "测试一轮",
			fields: fields{StructOfContainer: StructOfContainer{
				resolved: map[string]bool{"a": true},
				bindings: nil,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := &Container{
				StructOfContainer: tt.fields.StructOfContainer,
			}
			container.Flush()
		})
	}
}
