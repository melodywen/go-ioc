package container

import (
	"cjw.com/melodywen/go-ioc/mock"
	"fmt"
	"reflect"
	"testing"
)

func TestBodyOfContainer_DropStaleInstances(t *testing.T) {
	type fields struct {
		instances map[string]interface{}
	}
	type args struct {
		abstract interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name:   "测试如果是命中删除",
			fields: fields{instances: map[string]interface{}{"abc": 1}},
			args:   args{abstract: "abc"},
			want:   true,
		}, {
			name:   "测试如果是没有命中",
			fields: fields{instances: map[string]interface{}{"abc": 1}},
			args:   args{abstract: "abcc"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &BodyOfContainer{
				instances: tt.fields.instances,
			}
			if got := body.DropStaleInstances(tt.args.abstract); got != tt.want {
				t.Errorf("DropStaleInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBodyOfContainer_Resolved(t *testing.T) {
	type fields struct {
		instances map[string]interface{}
		resolved  map[string]bool
	}
	type args struct {
		abstract interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name: "如果在instance 中",
			fields: fields{
				instances: map[string]interface{}{"cjw.com/melodywen/go-ioc/mock.Animal": 1},
			},
			args: args{abstract: mock.Animal{}},
			want: true,
		},
		{
			name: "如果在resolved 中",
			fields: fields{
				resolved: map[string]bool{"cjw.com/melodywen/go-ioc/mock.Animal": true},
			},
			args: args{abstract: mock.Animal{}},
			want: true,
		},
		{
			name: "如果都不在中",
			fields: fields{
				resolved: map[string]bool{"*cjw.com/melodywen/go-ioc/mock.Animal": true},
			},
			args: args{abstract: mock.Animal{}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &BodyOfContainer{
				instances: tt.fields.instances,
				resolved:  tt.fields.resolved,
			}
			if got := body.Resolved(tt.args.abstract); got != tt.want {
				t.Errorf("Resolved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBodyOfContainer_Bind(t *testing.T) {
	type fields struct {
		instances         map[string]interface{}
		bindings          map[string]Bind
		resolved          map[string]bool
		BuildOfContainer  BuildOfContainer
		CommonOfContainer CommonOfContainer
	}
	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "测试concrete",
			fields: fields{
				bindings: map[string]Bind{"cjw.com/melodywen/go-ioc/mock.Animal": {
					shared:   true,
					concrete: 4,
				}},
			},
			args: args{
				abstract: mock.Animal{},
				concrete: 4,
				shared:   true,
			},
		},
		{
			name: "测试shared",
			fields: fields{
				bindings: map[string]Bind{"cjw.com/melodywen/go-ioc/mock.Animal": {
					shared:   false,
					concrete: mock.NewAnimal,
				}},
			},
			args: args{
				abstract: mock.Animal{},
				concrete: mock.NewAnimal,
				shared:   false,
			},
		},
		{
			name: "测试如果是绑定一个指针",
			fields: fields{
				bindings: map[string]Bind{"*cjw.com/melodywen/go-ioc/mock.Animal": {
					shared:   false,
					concrete: mock.NewAnimal,
				}},
			},
			args: args{
				abstract: &mock.Animal{},
				concrete: mock.NewAnimal,
				shared:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &BodyOfContainer{
				bindings: tt.fields.bindings,
			}
			obj := newBodyOfContainer()
			index := obj.AbstractToString(tt.args.abstract)
			obj.Bind(tt.args.abstract, tt.args.concrete, tt.args.shared)

			result := obj.bindings[index].shared != body.bindings[index].shared
			if reflect.TypeOf(obj.bindings[index].concrete).Kind() == reflect.Func {
				result = result || (reflect.ValueOf(obj.bindings[index].concrete).Pointer() != reflect.ValueOf(body.bindings[index].concrete).Pointer())
			} else {
				result = result || !reflect.DeepEqual(obj.bindings[index].concrete, body.bindings[index].concrete)
			}
			if result {
				fmt.Println(11111, obj.bindings[index])
				fmt.Println(11112, body.bindings[index])
				t.Errorf("Bind() = %v, want %v", body, obj)
			}
		})
	}
}

func TestBodyOfContainer_IsShared(t *testing.T) {
	type fields struct {
		instances         map[string]interface{}
		bindings          map[string]Bind
		resolved          map[string]bool
		BuildOfContainer  BuildOfContainer
		CommonOfContainer CommonOfContainer
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
		// TODO: Add test cases.
		{
			name:   "测试如果存在",
			fields: fields{bindings: map[string]Bind{"abc": {true, 1}}},
			args:   args{abstract: "abc"},
			want:   true,
		}, {
			name:   "测试如果存在",
			fields: fields{instances: map[string]interface{}{"abc": 1}},
			args:   args{abstract: "abc"},
			want:   true,
		}, {
			name:   "测试如果不在",
			fields: fields{bindings: map[string]Bind{"abc": {true, 1}}},
			args:   args{abstract: "abcc"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &BodyOfContainer{
				instances:         tt.fields.instances,
				bindings:          tt.fields.bindings,
				resolved:          tt.fields.resolved,
				BuildOfContainer:  tt.fields.BuildOfContainer,
				CommonOfContainer: tt.fields.CommonOfContainer,
			}
			if got := body.IsShared(tt.args.abstract); got != tt.want {
				t.Errorf("IsShared() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBodyOfContainer_GetConcrete(t *testing.T) {
	type fields struct {
		abstract interface{}
		concrete interface{}
		shared   bool
	}
	type args struct {
		abstract interface{}
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantConcrete interface{}
	}{
		// TODO: Add test cases.
		{
			name: "测试 数值 和 标量",
			fields: fields{
				abstract: 123,
				concrete: "hello",
			},
			args:         args{abstract: 123},
			wantConcrete: "hello",
		}, {
			name: "测试 struct 和 fun",
			fields: fields{
				abstract: mock.Animal{},
				concrete: mock.NewAnimal,
			},
			args:         args{abstract: mock.Animal{}},
			wantConcrete: mock.NewAnimal,
		}, {
			name: "测试 指正 和 fun",
			fields: fields{
				abstract: &mock.Animal{},
				concrete: mock.NewAnimal,
			},
			args:         args{abstract: &mock.Animal{}},
			wantConcrete: mock.NewAnimal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := newBodyOfContainer()
			body.Bind(tt.fields.abstract, tt.fields.concrete, true)
			gotConcrete := body.GetConcrete(tt.args.abstract)
			result := true
			if reflect.TypeOf(gotConcrete).Kind() == reflect.Func {
				result = reflect.ValueOf(gotConcrete).Pointer() != reflect.ValueOf(tt.wantConcrete).Pointer()
			} else {
				result = !reflect.DeepEqual(gotConcrete, tt.wantConcrete)
			}
			if result {
				t.Errorf("GetConcrete() = %v, want %v", gotConcrete, tt.wantConcrete)
			}
		})
	}
}
