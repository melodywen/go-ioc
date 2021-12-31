package container

import (
	"github.com/melodywen/go-ioc/mock"
	"reflect"
	"testing"
)

func TestContextualBindingBuilder_Give(t *testing.T) {
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
			name: "测试单值绑定",
			fields: fields{StructOfContainer: StructOfContainer{
				contextual: map[string]map[string]interface{}{
					"github.com/melodywen/go-ioc/mock.NewAnimalAndParam": {"github.com/melodywen/go-ioc/mock.Animal": 222},
				},
			}},
			args: args{
				concrete:       mock.NewAnimalAndParam,
				abstract:       mock.Animal{},
				implementation: 222,
			},
		}, {
			name: "测试多个值绑定",
			fields: fields{StructOfContainer: StructOfContainer{
				contextual: map[string]map[string]interface{}{
					"github.com/melodywen/go-ioc/mock.NewAnimalAndParam": {"github.com/melodywen/go-ioc/mock.Animal": 222},
					"github.com/melodywen/go-ioc/mock.NewAnimal":         {"github.com/melodywen/go-ioc/mock.Animal": 222},
				},
			}},
			args: args{
				concrete:       []interface{}{mock.NewAnimalAndParam, mock.NewAnimal},
				abstract:       mock.Animal{},
				implementation: 222,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewContainer()
			container.When(tt.args.concrete).Need(tt.args.abstract).Give(tt.args.implementation)
			if !reflect.DeepEqual(container.contextual, tt.fields.StructOfContainer.contextual) {
				t.Errorf("AddContextualBinding() = %v, want %v", container.StructOfContainer, tt.fields.StructOfContainer)
			}
		})
	}
}
