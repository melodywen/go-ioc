package container

import (
	"github.com/melodywen/go-ioc/mock"
	"fmt"
	"reflect"
	"testing"
)

func TestContainer_When(t *testing.T) {
	type fields struct {
		StructOfContainer StructOfContainer
		ExtendOfContainer ExtendOfContainer
	}
	container := NewContainer()
	type args struct {
		concrete interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ContextualBindingBuilder
	}{
		{
			name: "如果是单值",
			args: args{concrete: mock.Animal{}},
			want: newContextualBindingBuilder(container, []string{"github.com/melodywen/go-ioc/mock.Animal"}),
		}, {
			name: "如果是数组值",
			args: args{concrete: []interface{}{mock.Animal{}, mock.Animal{}}},
			want: newContextualBindingBuilder(container, []string{"github.com/melodywen/go-ioc/mock.Animal", "github.com/melodywen/go-ioc/mock.Animal"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := container.When(tt.args.concrete); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("When() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainer_AddContextualBinding(t *testing.T) {
	type fields struct {
		StructOfContainer StructOfContainer
	}
	type args struct {
		concrete       interface{}
		abstract       interface{}
		implementation interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "测试第一轮",
			fields: fields{StructOfContainer: StructOfContainer{
				contextual: map[string]map[string]interface{}{
					"github.com/melodywen/go-ioc/mock.Animal": {"github.com/melodywen/go-ioc/mock.AddAndParam": 222},
				},
			}},
			args: args{
				concrete:       mock.Animal{},
				abstract:       mock.AddAndParam,
				implementation: 222,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewContainer()
			container.AddContextualBinding(tt.args.concrete, tt.args.abstract, tt.args.implementation)

			fmt.Println(container.contextual)
			fmt.Println(tt.fields.StructOfContainer.contextual)
			if !reflect.DeepEqual(container.contextual, tt.fields.StructOfContainer.contextual) {
				t.Errorf("AddContextualBinding() = %v, want %v", container.StructOfContainer, tt.fields.StructOfContainer)
			}

		})
	}
}
