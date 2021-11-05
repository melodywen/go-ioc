package container

import (
	"cjw.com/melodywen/go-ioc/mock"
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
					concrete: mock.NewAnimal,
				}},
			},
			args: args{
				abstract: mock.Animal{},
				concrete: mock.NewAnimal,
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
			if obj.Bind(tt.args.abstract, tt.args.abstract, tt.args.shared); reflect.DeepEqual(obj, body) {
				t.Errorf("Bind() = %v, want %v", body, obj)
			}
		})
	}
}
