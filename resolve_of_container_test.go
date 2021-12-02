package container

import (
	"cjw.com/melodywen/go-ioc/mock"
	"fmt"
	"reflect"
	"testing"
)

func TestContainer_getConcrete(t *testing.T) {
	type fields struct {
		abstract interface{}
		concrete interface{}
		shared   bool
	}
	type args struct {
		abstract string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantConcrete interface{}
	}{
		{
			name: "第一轮测试,测试对象",
			fields: fields{
				abstract: mock.Animal{},
				concrete: mock.NewAnimal,
				shared:   false,
			},
			args:         args{abstract: "cjw.com/melodywen/go-ioc/mock.Animal"},
			wantConcrete: mock.NewAnimal,
		},
		{
			name: "第一轮测试,如果不存在呢",
			fields: fields{
				abstract: &mock.Animal{},
				concrete: mock.NewAnimal,
				shared:   false,
			},
			args:         args{abstract: "cjw.com/melodywen/go-ioc/mock.Animal"},
			wantConcrete: "cjw.com/melodywen/go-ioc/mock.Animal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewContainer()
			container.Bind(tt.fields.abstract, tt.fields.concrete, tt.fields.shared)
			gotConcrete := container.getConcrete(tt.args.abstract)

			if reflect.TypeOf(tt.wantConcrete).Kind() == reflect.Func {
				sf1 := reflect.ValueOf(gotConcrete)
				sf2 := reflect.ValueOf(tt.wantConcrete)
				if sf1.Pointer() != sf2.Pointer() {
					t.Errorf("getConcrete() = %v, want %v", gotConcrete, tt.wantConcrete)
				}
			} else {
				if !reflect.DeepEqual(gotConcrete, tt.wantConcrete) {
					t.Errorf("getConcrete() = %v, want %v", gotConcrete, tt.wantConcrete)
				}
			}

		})
	}
}

func TestContainer_IsShared(t *testing.T) {
	type fields struct {
		abstract interface{}
		concrete interface{}
		shared   bool
		instance interface{}
	}
	type args struct {
		abstract string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "测试——非分享",
			fields: fields{
				abstract: mock.NewAnimal,
				concrete: 1,
				shared:   false,
			},
			args: args{abstract: "cjw.com/melodywen/go-ioc/mock.NewAnimal"},
			want: false,
		},
		{
			name: "测试——分享",
			fields: fields{
				abstract: mock.NewAnimal,
				concrete: 1,
				shared:   true,
			},
			args: args{abstract: "cjw.com/melodywen/go-ioc/mock.NewAnimal"},
			want: true,
		}, {
			name: "测试——不存在",
			fields: fields{
				abstract: mock.NewAnimal,
				concrete: 1,
				shared:   true,
			},
			args: args{abstract: "abcc"},
			want: false,
		}, {
			name: "测试——如果在缓存中",
			fields: fields{
				abstract: mock.NewAnimal,
				concrete: 1,
				shared:   true,
				instance: 1,
			},
			args: args{abstract: "cjw.com/melodywen/go-ioc/mock.NewAnimal"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewContainer()
			if tt.fields.instance != nil {
				container.Instance(tt.fields.abstract, tt.fields.instance)
			} else {
				container.Bind(tt.fields.abstract, tt.fields.concrete, tt.fields.shared)
			}
			if got := container.IsShared(tt.args.abstract); got != tt.want {
				fmt.Println(container)
				t.Errorf("IsShared() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_MakeWithParams(t *testing.T) {
	type fields struct {
		abstract interface{}
		concrete interface{}
		shared   bool
		alias    string
	}
	type args struct {
		abstract   interface{}
		parameters []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "测试一个加法",
			fields: fields{
				abstract: mock.AddNum,
				concrete: mock.AddNum,
				shared:   false,
			},
			args: args{
				abstract:   mock.AddNum,
				parameters: []interface{}{1, 2},
			},
			want: 3,
		}, {
			name: "测试如果是实例化对象",
			fields: fields{
				abstract: "abc",
				concrete: mock.NewAnimal,
				shared:   true,
			},
			args: args{
				abstract:   "abc",
				parameters: []interface{}{"dog", 2, "cate"},
			},
			want: mock.NewAnimal("dog", 2, "cate"),
		}, {
			name: "测试如果是实例化对象是个复杂的函数",
			fields: fields{
				abstract: "many-return",
				concrete: mock.NewAnimalAndParam,
				shared:   true,
			},
			args: args{
				abstract:   "many-return",
				parameters: []interface{}{"dog", 2, "cate"},
			},
			want: []interface{}{mock.NewAnimal("dog", 2, "cate"), "dog", 2, "cate"},
		}, {
			name: "测试如果是实例化对象是个复杂的函数 -> 设置缓存",
			fields: fields{
				abstract: "many-return-set",
				concrete: func() (*mock.Animal, string, int, string) {
					return mock.NewAnimalAndParam("dog", 2, "cate")
				},
				shared: true,
			},
			args: args{
				abstract:   "many-return-set",
				parameters: []interface{}{},
			},
			want: []interface{}{mock.NewAnimal("dog", 2, "cate"), "dog", 2, "cate"},
		}, {
			name: "测试如果是实例化对象是个复杂的函数- 获取有缓存",
			fields: fields{
				abstract: nil,
				concrete: mock.NewAnimalAndParam,
				shared:   true,
			},
			args: args{
				abstract:   "many-return-set",
				parameters: []interface{}{},
			},
			want: []interface{}{mock.NewAnimal("dog", 2, "cate"), "dog", 2, "cate"},
		}, {
			name: "测试如果是实例化对象是个复杂的函数- bind 重复绑定",
			fields: fields{
				abstract: "alias-many-return-set-alias",
				concrete: "many-return-set",
				shared:   true,
			},
			args: args{
				abstract:   "alias-many-return-set-alias",
				parameters: []interface{}{},
			},
			want: []interface{}{mock.NewAnimal("dog", 2, "cate"), "dog", 2, "cate"},
		},
	}
	container := NewContainer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.abstract != nil {
				container.Bind(tt.fields.abstract, tt.fields.concrete, tt.fields.shared)
			}
			fmt.Println(container)
			got := container.MakeWithParams(tt.args.abstract, tt.args.parameters)
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Println(got, tt.want)
				t.Errorf("MakeWithParams() = %v, want %v", got, tt.want)
			} else {
				//fmt.Println(container)
			}
		})
	}
}

func TestContainer_makeWithBuildStack(t *testing.T) {
	type args struct {
		abstract   interface{}
		shared   bool
		concrete interface{}
		parameters []interface{}
		buildStack []string
	}
	tests := []struct {
		name   string
		args   args
		want   interface{}
	}{
		{
			name:"测试简单的stack 注册",
			args: args{
				abstract: mock.Animal{},
				shared: true,
				concrete: mock.NewAnimal,
				parameters: []interface{}{"cat",3,"cate"},
				buildStack: nil,
			},
			want: mock.NewAnimal("cat",3,"cate"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewContainer()
			container.Bind(tt.args.abstract,tt.args.concrete,tt.args.shared)
			got := container.MakeWithParams(tt.args.abstract, tt.args.parameters);
			fmt.Println(got)
			if  !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeWithBuildStack() = %v, want %v", got, tt.want)
			}
		})
	}
}