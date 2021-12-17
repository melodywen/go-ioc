package container

import (
	"github.com/melodywen/go-ioc/mock"
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

func TestNewFamily(t *testing.T) {
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
			name: "register-work",
			args: args{
				abstract: mock.Work{},
				concrete: func() *mock.Work { return mock.NewWork("php工程师", 12) },
				shared:   false,
				param:    []interface{}{},
			},
			want: mock.NewWork("php工程师", 12),
		}, {
			name: "register-homeWork",
			args: args{
				abstract: mock.Homework{},
				concrete: func() *mock.Homework { return mock.NewHomework("照顾家庭") },
				shared:   false,
				param:    []interface{}{},
			},
			want: mock.NewHomework("照顾家庭"),
		}, {
			name: "register-father",
			args: args{
				abstract: mock.Father{},
				concrete: func(work *mock.Work) *mock.Father { return mock.NewFatherWithAllParam(work, "张三", 33) },
				shared:   false,
				param:    []interface{}{},
			},
			want: mock.NewFatherWithAllParam(mock.NewWork("php工程师", 12), "张三", 33),
		}, {
			name: "register-mother",
			args: args{
				abstract: mock.Mother{},
				concrete: func(homework mock.Homework) *mock.Mother { return mock.NewMother("赵丽", 18, homework) },
				shared:   false,
				param:    []interface{}{},
			},
			want: mock.NewMother("赵丽", 18, *mock.NewHomework("照顾家庭")),
		}, {
			name: "register-family",
			args: args{
				abstract: mock.Family{},
				concrete: mock.NewFamily,
				shared:   false,
				param: []interface{}{
					"张三的幸福家庭",
					[]mock.Child{
						*mock.NewChild("熊大", 12),
						*mock.NewChild("熊二", 10),
					},
				},
			},
			want: mock.NewFamily(
				"张三的幸福家庭",
				*mock.NewFatherWithAllParam(mock.NewWork("php工程师", 12), "张三", 33),
				[]mock.Child{
					*mock.NewChild("熊大", 12),
					*mock.NewChild("熊二", 10),
				},
				*mock.NewMother("赵丽", 18, *mock.NewHomework("照顾家庭")),
			),
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

func TestNewHomework(t *testing.T) {
	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
		param    []interface{}
	}
	var workInterface *mock.WorkInterface

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "register-work-interface",
			args: args{
				abstract: workInterface,
				concrete: func() interface{} { return nil },
				shared:   false,
				param:    []interface{}{},
			},
			want: nil,
		}, {
			name: "register-work",
			args: args{
				abstract: mock.Work{},
				concrete: func() *mock.Work { return mock.NewWork("php工程师", 12) },
				shared:   false,
				param:    []interface{}{},
			},
			want: mock.NewWork("php工程师", 12),
		},
		{
			name: "register-homeWork",
			args: args{
				abstract: mock.Homework{},
				concrete: func() *mock.Homework { return mock.NewHomework("照顾家庭") },
				shared:   false,
				param:    []interface{}{},
			},
			want: mock.NewHomework("照顾家庭"),
		},
		{
			name: "register-father",
			args: args{
				abstract: mock.Father{},
				concrete: mock.NewFatherWithInterface,
				shared:   false,
				param:    []interface{}{"张三", 33},
			},
			want: mock.NewFatherWithAllParam(mock.NewWork("php工程师", 12), "张三", 33),
		},
		{
			name: "register-mother",
			args: args{
				abstract: mock.Mother{},
				concrete: mock.NewMotherWithInterface,
				shared:   false,
				param:    []interface{}{"赵丽", 18},
			},
			want: mock.NewMother("赵丽", 18, *mock.NewHomework("照顾家庭")),
		},
	}
	app := NewContainer()
	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if index == 2 {
				app.When(mock.NewFatherWithInterface).Need(workInterface).Give(func(work *mock.Work) *mock.Work {
					return work
				})
				app.When(mock.NewMotherWithInterface).Need(workInterface).Give(func(work mock.Homework) mock.Homework {
					return work
				})
			}
			app.Bind(tt.args.abstract, tt.args.concrete, tt.args.shared)
			got := app.MakeWithParams(tt.args.abstract, tt.args.param)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWork() = %v, want %v", got, tt.want)
			}
		})
	}
}
