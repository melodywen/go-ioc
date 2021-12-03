package container

import (
	"cjw.com/melodywen/go-ioc/mock"
	"reflect"
	"testing"
)

func TestNewWork(t *testing.T) {
	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
		param    []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "简单work",
			args: args{
				abstract: mock.Work{},
				concrete: mock.NewWork,
				shared:   true,
				param:    []interface{}{"php工程师", 12},
			},
			want: mock.NewWork("php工程师", 12),
		}, {
			name: "简单homeWork",
			args: args{
				abstract: mock.Homework{},
				concrete: mock.NewHomework,
				shared:   true,
				param:    []interface{}{"照顾家庭"},
			},
			want: mock.NewHomework("照顾家庭"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewContainer()
			app.Bind(tt.args.abstract, tt.args.concrete, tt.args.shared)
			got := app.MakeWithParams(tt.args.abstract, tt.args.param)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWork() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFather(t *testing.T) {
	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
		param    []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "简单work",
			args: args{
				abstract: mock.Work{},
				concrete: func() *mock.Work { return mock.NewWork("php工程师", 12) },
				shared:   false,
				param:    []interface{}{},
			},
			want: mock.NewWork("php工程师", 12),
		}, {
			name: "简单homeWork",
			args: args{
				abstract: mock.Homework{},
				concrete: mock.NewHomework,
				shared:   true,
				param:    []interface{}{"照顾家庭"},
			},
			want: mock.NewHomework("照顾家庭"),
		}, {
			name: "简单father-如果是所有的参数",
			args: args{
				abstract: mock.Father{},
				concrete: mock.NewFatherWithAllParam,
				shared:   true,
				param:    []interface{}{"张三", 33},
			},
			want: mock.NewFatherWithAllParam(mock.NewWork("php工程师", 12), "张三", 33),
		}, {
			name: "简单father-如果是struct单个参数",
			args: args{
				abstract: mock.Father{},
				concrete: mock.NewFatherWithStruct,
				shared:   true,
				param:    []interface{}{"张三", 33},
			},
			want: mock.NewFatherWithAllParam(mock.NewWork("php工程师", 12), "张三", 33),
		}, {
			name: "简单father-如果是pre单个参数",
			args: args{
				abstract: mock.Father{},
				concrete: mock.NewFatherWithPre,
				shared:   true,
				param:    []interface{}{"张三", 33},
			},
			want: mock.NewFatherWithAllParam(mock.NewWork("php工程师", 12), "张三", 33),
		},
	}
	app := NewContainer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.Bind(tt.args.abstract, tt.args.concrete, tt.args.shared)
			got := app.MakeWithParams(tt.args.abstract, tt.args.param)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWork() = %v, want %v", got, tt.want)
			}
		})
	}
}
