package container

import (
	"testing"
)

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
