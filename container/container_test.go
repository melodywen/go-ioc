package container

import (
	"fmt"
	"github.com/melodywen/supports/exceptions"
	"reflect"
	"testing"
)

// containerChild
//  @Description:
type containerChild struct {
	Container
}

func newContainerChild() *containerChild {
	child := &containerChild{}
	child.Container = *NewContainer()
	child.child = child
	return child
}

// ResolveCallback
//  @Description:
//  @receiver c
//  @param identifier
func (c *containerChild) ResolveCallback(identifier string) {
	//fmt.Println("callback child ~~~~~")
}

var container *containerChild

func errRecover() {
	e := recover()
	if e == nil {
		return
	}
	fmt.Println(e)
	err := e.(exceptions.ErrorInterface)
	fmt.Println("错误详情", err.ErrorInspect())
}

func init() {
	container = newContainerChild()
}

func TestNewContainer(t *testing.T) {
	t.Run("创建container", func(t *testing.T) {
		container = newContainerChild()
		fmt.Println("步骤-1: 创建 子容器 - ", container)
	})
}

func TestContainer_Build(t *testing.T) {
	container.Flush()
	container.SingletonIf(&Request{}, newRequest)
	container.Instance(Request{}, *newRequest())
	container.SingletonIf(&Response{}, newResponse)
	container.Instance(Response{}, *newResponse())

	type args struct {
		concrete   any
		parameters []any
		stack      *containerStack
	}
	tests := []struct {
		name       string
		args       args
		wantObject any
	}{
		{
			name: "测试返回单个值",
			args: args{
				concrete:   newUserController,
				parameters: nil,
				stack:      &containerStack{},
			},
			wantObject: newUserController(),
		}, {
			name: "测试返回多个值",
			args: args{
				concrete:   newUserControllerAndOther,
				parameters: []any{newRequest(), *newResponse(), "张三", 12},
				stack:      &containerStack{},
			},
			wantObject: []any{newUserController(), *newRequest(), newResponse(), "张三", 12},
		}, {
			name: "测试自动注入",
			args: args{
				concrete:   newUserControllerAndObj,
				parameters: nil,
				stack:      &containerStack{},
			},
			wantObject: []any{newUserController(), *newRequest(), newResponse()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotObject := container.Build(tt.args.concrete, tt.args.parameters, tt.args.stack)
			if !reflect.DeepEqual(gotObject, tt.wantObject) {
				t.Errorf("Build() = %v, want %v", gotObject, tt.wantObject)
			}
		})
	}
}
func TestContainer_Build_ERR(t *testing.T) {
	container.Flush()
	container.SingletonIf(&Request{}, newRequest)
	container.Instance(Request{}, *newRequest())
	container.SingletonIf(&Response{}, newResponse)
	container.Instance(Response{}, *newResponse())
	type args struct {
		concrete   any
		parameters []any
		stack      *containerStack
	}
	tests := []struct {
		name       string
		args       args
		wantObject any
	}{
		{
			name: "自动加载实参类型错误",
			args: args{
				concrete:   newUserControllerAndOther,
				parameters: nil,
				stack:      &containerStack{},
			},
			wantObject: nil,
		}, {
			name: "如果不是函数",
			args: args{
				concrete:   newUserController(),
				parameters: nil,
				stack:      &containerStack{},
			},
			wantObject: nil,
		}, {
			name: "自动加载参数未注入",
			args: args{
				concrete:   newGorm,
				parameters: nil,
				stack:      &containerStack{},
			},
			wantObject: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer errRecover()
			gotObject := container.Build(tt.args.concrete, tt.args.parameters, tt.args.stack)
			if !reflect.DeepEqual(gotObject, tt.wantObject) {
				t.Errorf("Build() = %v, want %v", gotObject, tt.wantObject)
			}
		})
	}
}

func TestContainer_Callback(t *testing.T) {
	container.Flush()
	container.SingletonIf(&Request{}, newRequest)
	container.Instance(Request{}, *newRequest())
	container.SingletonIf(&Response{}, newResponse)
	container.Instance(Response{}, *newResponse())

	container.BindIf(newUserControllerAndObj, nil, false)

	container.BeforeResolving(nil, func(s string, param []any, c *Container) {
		fmt.Println("global-before-callback", s)
	})
	container.BeforeResolving(newUserControllerAndObj, func(s string, param []any, c *Container) {
		fmt.Println("before-callback", s)
	})
	container.Resolving(nil, func(instance any, container *Container) {
		fmt.Println("global-ing-callback", instance)
	})
	container.Resolving(newUserControllerAndObj, func(instance any, container *Container) {
		fmt.Println("ing-callback", instance)
	})
	container.AfterResolving(nil, func(instance any, container *Container) {
		fmt.Println("global-after-callback", instance)
	})
	container.AfterResolving(newUserControllerAndObj, func(instance any, container *Container) {
		fmt.Println("after-callback", instance)
	})

	t.Run("查询回调-callback", func(t *testing.T) {
		gotObject := container.Make(newUserControllerAndObj)
		s1, s2, s3 := newUserControllerAndObj(newRequest(), *newResponse())
		wantObject := []any{s1, s2, s3}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})

	container.Extend(newUserControllerAndObj, func(object any, container *Container) any {
		m := object.([]any)
		response := m[2].(*Response)
		response.Msg += "callback1"
		return object
	})
	container.Extend(newUserControllerAndObj, func(object any, container *Container) any {
		m := object.([]any)
		response := m[2].(*Response)
		response.Msg += "callback2"
		return object
	})
	t.Run("查询回调-extend1", func(t *testing.T) {
		gotObject := container.Make(newUserControllerAndObj)
		s1, s2, s3 := newUserControllerAndObj(newRequest(), *newResponse())
		s3.Msg += "callback1callback2"
		wantObject := []any{s1, s2, s3}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})

	container.Bind(newUserControllerAndObj, nil, true)
	container.Make(newUserControllerAndObj)
	container.Extend(newUserControllerAndObj, func(object any, container *Container) any {
		m := object.([]any)
		response := m[2].(*Response)
		response.Msg += "callback3"
		return object
	})

	container.Make(newUserControllerAndObj)
	t.Run("查询回调-extend2", func(t *testing.T) {
		gotObject := container.Make(newUserControllerAndObj)
		s1, s2, s3 := newUserControllerAndObj(newRequest(), *newResponse())
		s3.Msg += "callback1callback2callback3"
		wantObject := []any{s1, s2, s3}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})
	container.ForgetExtenders(newUserControllerAndObj)

	container.Rebinding(newUserControllerAndObj, func(container *Container, instance any) {
		fmt.Println("rebind callback", AbstractToString(newUserControllerAndObj))
	})
	container.Rebinding(newUserControllerAndOther, func(container *Container, instance any) {
		fmt.Println("rebind callback", AbstractToString(newUserControllerAndOther))
	})
	container.Bind(newUserControllerAndObj, nil, false)
}

func TestContainer_When(t *testing.T) {
	container.Flush()

	container.Bind(&ossAli{}, newOssAli, true)
	container.Bind(&ossTencent{}, newOssTencent, true)

	container.Bind(newUserControllerAndOss, nil, true)
	container.Bind(newFileControllerAndOss, nil, true)

	var o *ossInterface
	container.When([]any{newUserControllerAndOss}).Need(o).Give(func(oss *ossAli) ossInterface {
		return oss
	})
	container.When([]any{newFileControllerAndOss}).Need(o).Give(func(oss *ossTencent) ossInterface {
		return oss
	})
	t.Run("上下文查询-user controller", func(t *testing.T) {
		gotObject := container.Make(newUserControllerAndOss)
		wantObject := []any{newUserController(), newOssAli()}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})

	t.Run("上下文查询-file controller", func(t *testing.T) {
		defer errRecover()
		gotObject := container.Make(newFileControllerAndOss)
		wantObject := []any{newFileController(), newOssTencent()}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})

	// 为oss 起别名
	container.Flush()
	container.Bind(&ossAli{}, newOssAli, true)
	container.Bind(&ossTencent{}, newOssTencent, true)
	container.Bind(newUserControllerAndOss, nil, true)
	container.Bind(newFileControllerAndOss, nil, true)
	container.When([]any{newUserControllerAndOss}).Need("oss").Give(func(oss *ossAli) ossInterface {
		return oss
	})
	container.When([]any{newFileControllerAndOss}).Need("oss").Give(func(oss *ossTencent) ossInterface {
		return oss
	})
	container.Alias(o, "oss")
	t.Run("上下文查询-绑定别名-user controller", func(t *testing.T) {
		gotObject := container.Make(newUserControllerAndOss)
		wantObject := []any{newUserController(), newOssAli()}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
		wantObject[1].(ossInterface).Connect()
	})
	t.Run("上下文查询-绑定别名-file controller", func(t *testing.T) {
		defer errRecover()
		gotObject := container.Make(newFileControllerAndOss)
		wantObject := []any{newFileController(), newOssTencent()}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
		wantObject[1].(ossInterface).Connect()
	})

	// 为oss 起别名 但是没有命中
	container.Flush()
	container.Bind(&ossAli{}, newOssAli, true)
	container.Bind(&ossTencent{}, newOssTencent, true)
	container.Bind(newUserControllerAndOss, nil, true)
	container.Bind(newFileControllerAndOss, nil, true)
	container.When([]any{newUserControllerAndOss}).Need("oss").Give(func(oss *ossAli) ossInterface {
		return oss
	})
	container.When([]any{newFileControllerAndOss}).Need("oss").Give(func(oss *ossTencent) ossInterface {
		return oss
	})
	container.Alias(o, "oss-normal")
	container.Bind(o, newOss, false)

	t.Run("上下文查询-绑定别名-user controller", func(t *testing.T) {
		defer errRecover()
		gotObject := container.Make(newUserControllerAndOss)
		wantObject := []any{newUserController(), newOss()}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
		wantObject[1].(ossInterface).Connect()
	})
	t.Run("上下文查询-绑定别名-file controller", func(t *testing.T) {
		defer errRecover()
		gotObject := container.Make(newFileControllerAndOss)
		wantObject := []any{newFileController(), newOss()}
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})
}

func TestAbstractToStringErr(t *testing.T) {
	type args struct {
		abstract any
	}
	tests := []struct {
		name           string
		args           args
		wantIdentifier string
	}{
		{
			name: "如果是数字",
			args: args{abstract: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer errRecover()
			if gotIdentifier := AbstractToString(tt.args.abstract); gotIdentifier != tt.wantIdentifier {
				t.Errorf("AbstractToString() = %v, want %v", gotIdentifier, tt.wantIdentifier)
			}
		})
	}
}

func TestContainer_Factory(t *testing.T) {
	container.Flush()
	container.SingletonIf(&mysql{}, newMysql)
	type args struct {
		abstract any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "factory mysql",
			args: args{abstract: &mysql{}},
			want: newMysql(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFun := container.Factory(tt.args.abstract)
			got := gotFun()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Factory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_MakeWithParams(t *testing.T) {
	container.Flush()
	container.SingletonIf(&gorm{}, newGorm)
	type args struct {
		abstract   any
		parameters []any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "依赖参数",
			args: args{
				abstract:   &gorm{},
				parameters: []any{*newMysql()},
			},
			want: newGorm(*newMysql()),
		}, {
			name: "依赖参数",
			args: args{
				abstract:   &gorm{},
				parameters: []any{*newMysql()},
			},
			want: newGorm(*newMysql()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := container.MakeWithParams(tt.args.abstract, tt.args.parameters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeWithParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestContainer_resolve(t *testing.T) {
	container.Flush()
	container.Instance("mysql", *newMysql())
	container.Singleton(mysql{}, "mysql")
	container.Singleton(newGorm, nil)

	t.Run("resolve - 通过构造回调函数", func(t *testing.T) {
		gotObject := container.Make(newGorm)
		wantObject := newGorm(*newMysql())
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})

	t.Run("resolve - 不存在", func(t *testing.T) {
		defer errRecover()
		gotObject := container.Make("none")
		wantObject := newGorm(*newMysql())
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})
	container.Bind("none", nil, false)
	t.Run("resolve - 不存在", func(t *testing.T) {
		defer errRecover()
		gotObject := container.Make("none")
		wantObject := newGorm(*newMysql())
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})
}

func TestContainer_Store(t *testing.T) {
	container.Flush()
	container.Bind("mysql", newMysql, true)
	container.Alias(mysql{}, "mysql")
	container.Instance("mysql", *newMysql())

	t.Run("检测 - Instance", func(t *testing.T) {
		gotObject := container.Make("mysql")
		wantObject := *newMysql()
		if !reflect.DeepEqual(gotObject, wantObject) {
			t.Errorf("Build() = %v, want %v", gotObject, wantObject)
		}
	})

	t.Run("检测 - Alias", func(t *testing.T) {
		defer errRecover()
		container.Alias("a", "a")
	})

	t.Run("检测 - Resolved", func(t *testing.T) {
		container.Resolved("mysql")
		container.Resolved(mysql{})
	})
	container.Alias(mysql{}, "mysql")
	container.Alias("mysql", "mysql2")
	container.Alias("mysql2", "mysql3")
	t.Run("检测 - alias", func(t *testing.T) {
		container.Resolved("mysql3")
	})
}
