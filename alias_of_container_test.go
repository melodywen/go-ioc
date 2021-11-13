package container

import (
	"cjw.com/melodywen/go-ioc/mock"
	"reflect"
	"testing"
)

func TestContainer_IsAlias(t *testing.T) {
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
			name: "测试true",
			fields: fields{StructOfContainer{
				aliases: map[string]string{"aa": "aa"},
			}},
			args:   args{abstract: "aa"},
			wantOk: true,
		}, {
			name: "测试false",
			fields: fields{StructOfContainer{
				aliases: map[string]string{"aa": "aa"},
			}},
			args:   args{abstract: "ab"},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := &Container{
				StructOfContainer: tt.fields.StructOfContainer,
			}
			if gotOk := container.IsAlias(tt.args.abstract); gotOk != tt.wantOk {
				t.Errorf("IsAlias() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestContainer_GetAlias(t *testing.T) {
	type fields struct {
		StructOfContainer StructOfContainer
	}
	type args struct {
		abstract interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "测试返回自己，并且递归",
			fields: fields{StructOfContainer{
				aliases: map[string]string{"aa": "aabbcc", "aabbcc": "abc"},
			}},
			args: args{abstract: "aa"},
			want: "abc",
		}, {
			name: "测试返回自己",
			fields: fields{StructOfContainer{
				aliases: map[string]string{"aa": "aabbcc"},
			}},
			args: args{abstract: "abc"},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := &Container{
				StructOfContainer: tt.fields.StructOfContainer,
			}
			if got := container.GetAlias(tt.args.abstract); got != tt.want {
				t.Errorf("GetAlias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_Alias(t *testing.T) {
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
	}{
		{
			name: "测试一组",
			fields: fields{StructOfContainer{
				aliases:         map[string]string{"dog": "cjw.com/melodywen/go-ioc/mock.NewAnimal"},
				abstractAliases: map[string][]string{"cjw.com/melodywen/go-ioc/mock.NewAnimal": []string{"dog"}},
			}},
			args: args{
				abstract: mock.NewAnimal,
				alias:    "dog",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := newContainer()
			container.Alias(tt.args.abstract, tt.args.alias)
			if !reflect.DeepEqual(tt.fields.StructOfContainer.aliases, container.aliases) ||
				!reflect.DeepEqual(tt.fields.StructOfContainer.abstractAliases, container.abstractAliases) {
				t.Errorf("GetAlias() = %v, want %v", tt.fields, container)
			}

		})
	}
}
