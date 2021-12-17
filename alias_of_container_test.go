package container

import (
	"github.com/melodywen/go-ioc/mock"
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
				aliases:         map[string]string{"dog": "github.com/melodywen/go-ioc/mock.NewAnimal"},
				abstractAliases: map[string][]string{"github.com/melodywen/go-ioc/mock.NewAnimal": []string{"dog"}},
			}},
			args: args{
				abstract: mock.NewAnimal,
				alias:    "dog",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewContainer()
			container.Alias(tt.args.abstract, tt.args.alias)
			if !reflect.DeepEqual(tt.fields.StructOfContainer.aliases, container.aliases) ||
				!reflect.DeepEqual(tt.fields.StructOfContainer.abstractAliases, container.abstractAliases) {
				t.Errorf("GetAlias() = %v, want %v", tt.fields, container)
			}

		})
	}
}

func TestContainer_removeAbstractAlias(t *testing.T) {
	type fields struct {
		abstract interface{}
		alias    interface{}
	}
	type args struct {
		search string
	}
	tests := []struct {
		name   string
		fields []fields
		args   args
		wantOk bool
	}{
		{
			name: "简单测试单个",
			fields: []fields{
				{
					alias:    "abc",
					abstract: mock.NewAnimalAndParam,
				},
			},
			args:   args{search: "abc"},
			wantOk: true,
		},
		{
			name: "简单测试多个",
			fields: []fields{
				{
					alias:    "bbc",
					abstract: mock.NewAnimalAndParam,
				}, {
					alias:    "bbcc",
					abstract: mock.NewAnimalAndParam,
				},
			},
			args:   args{search: "abc"},
			wantOk: false,
		},
	}
	container := NewContainer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, field := range tt.fields {
				container.Alias(field.abstract, field.alias)
			}
			gotOk := container.removeAbstractAlias(tt.args.search)
			if gotOk != tt.wantOk {
				t.Errorf("removeAbstractAlias() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
