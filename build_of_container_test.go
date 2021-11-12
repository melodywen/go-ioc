package container

import (
	"cjw.com/melodywen/go-ioc/mock"
	"reflect"
	"testing"
)

func TestBuildOfContainer_Build(t *testing.T) {
	type args struct {
		concrete   interface{}
		parameters []interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantObject interface{}
	}{
		// TODO: Add test cases.
		{
			name: "测试如果是一个回调函数，并且返回一个参数",
			args: args{
				concrete:   mock.AddNum,
				parameters: []interface{}{1, 2},
			},
			wantObject: 3,
		},
		{
			name: "测试如果是一个回调函数，并且返回多个参数",
			args: args{
				concrete:   mock.AddAndParam,
				parameters: []interface{}{1, 2},
			},
			wantObject: []interface{}{3, 1, 2},
		},
		{
			name: "测试如果是一个回调函数，实例化对象",
			args: args{
				concrete:   mock.NewAnimal,
				parameters: []interface{}{"小猫", 2, "猫科"},
			},
			wantObject: mock.NewAnimal("小猫", 2, "猫科"),
		},
		{
			name: "测试如果是一个回调函数，并且返回多个参数",
			args: args{
				concrete:   mock.NewAnimalAndParam,
				parameters: []interface{}{"小猫", 2, "猫科"},
			},
			wantObject: []interface{}{
				mock.NewAnimal("小猫", 2, "猫科"),
				"小猫", 2, "猫科",
			},
		},
		{
			name: "测试如果是一个标量,则直接报错",
			args: args{
				concrete:   []interface{}{"小猫", 2, "猫科"},
				parameters: []interface{}{"小猫", 2, "猫科"},
			},
			wantObject: []interface{}{
				mock.NewAnimal("小猫", 2, "猫科"),
				"小猫", 2, "猫科",
			},
		},
		// todo 待实现 （完成动态代理功能）
		//{
		//	name: "测试如果是一个结构体，则直接进行实例化",
		//	args: args{
		//		concrete:   mock.Animal{},
		//		parameters: []interface{}{"小猫", 2, "猫科"},
		//	},
		//	wantObject: []interface{}{
		//		mock.NewAnimal("小猫", 2, "猫科"),
		//		"小猫", 2, "猫科",
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bu := &BuildOfContainer{}
			//gotObject := bu.Build(tt.args.concrete, tt.args.parameters)
			//fmt.Println(gotObject,reflect.TypeOf(gotObject))
			if gotObject := bu.Build(tt.args.concrete, tt.args.parameters); !reflect.DeepEqual(gotObject, tt.wantObject) {
				t.Errorf("Build() = %v, want %v", gotObject, tt.wantObject)
			}
		})
	}
}

func TestBuildOfContainer_isBuildable(t *testing.T) {
	type args struct {
		abstract interface{}
		concrete interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "简单测试一下",
			args: args{
				abstract: mock.Animal{},
				concrete: mock.NewAnimal,
			},
			want: true,
		}, {
			name: "如果不是回调false",
			args: args{
				abstract: mock.Animal{},
				concrete: 3,
			},
			want: false,
		}, {
			name: "如果完全雷同",
			args: args{
				abstract: mock.Animal{},
				concrete: mock.Animal{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bu := &BuildOfContainer{}
			if got := bu.isBuildable(tt.args.abstract, tt.args.concrete); got != tt.want {
				t.Errorf("isBuildable() = %v, want %v", got, tt.want)
			}
		})
	}
}
