package container

import (
	"reflect"
	"testing"
)

func BenchmarkContainer_Build(t *testing.B) {
	type args struct {
		concrete   any
		parameters []any
		stack      *containerStack
	}
	tests := []struct {
		name       string
		args       args
		wantObject interface{}
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
			name: "测试返回单个值-不用反射",
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
			wantObject: []any{
				newUserController(), *newRequest(), newResponse(), "张三", 12,
			},
		}, {
			name: "测试返回多个值-不用反射",
			args: args{
				concrete:   newUserControllerAndOther,
				parameters: []any{newRequest(), *newResponse(), "张三", 12},
				stack:      &containerStack{},
			},
			wantObject: []any{
				newUserController(), *newRequest(), newResponse(), "张三", 12,
			},
		},
	}
	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			container := NewContainer()
			var gotObject any
			if index == 1 {
				gotObject = newUserController()
			} else if index == 3 {
				r1, r2, r3, r4, r5 := newUserControllerAndOther(newRequest(), *newResponse(), "张三", 12)
				gotObject = []any{r1, r2, r3, r4, r5}
			} else {
				gotObject = container.Build(tt.args.concrete, tt.args.parameters, tt.args.stack)
			}

			if !reflect.DeepEqual(gotObject, tt.wantObject) {
				t.Errorf("Build() = %v, want %v", gotObject, tt.wantObject)
			}
		})
	}
}
