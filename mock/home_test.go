package mock

import (
	container "cjw.com/melodywen/go-ioc"
	"reflect"
	"testing"
)

func TestNewWork(t *testing.T) {
	type args struct {
		abstract interface{}
		concrete interface{}
		shared   bool
		param []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "简单work",
			args: args{
				abstract: Work{},
				concrete: NewWork,
				shared:   true,
				param:  []interface{}{"php工程师", 12},
			},
			want: NewWork("php工程师", 12),
		},{
			name: "简单homeWork",
			args: args{
				abstract: Homework{},
				concrete: NewHomework,
				shared:   true,
				param:  []interface{}{"照顾家庭"},
			},
			want: NewHomework("照顾家庭"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := container.NewContainer()
			app.Bind(tt.args.abstract, tt.args.concrete, tt.args.shared)
			got := app.MakeWithParams(tt.args.abstract,tt.args.param)
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
		param []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "简单work",
			args: args{
				abstract: Work{},
				concrete: NewWork,
				shared:   true,
				param:  []interface{}{"php工程师", 12},
			},
			want: NewWork("php工程师", 12),
		},{
			name: "简单homeWork",
			args: args{
				abstract: Homework{},
				concrete: NewHomework,
				shared:   true,
				param:  []interface{}{"照顾家庭"},
			},
			want: NewHomework("照顾家庭"),
		},{
			name: "简单father",
			args: args{
				abstract: Father{},
				concrete: NewFather,
				shared:   true,
				param:  []interface{}{"张三",33},
			},
			want: NewFather("张三",33, NewWork("php工程师", 12)),
		},
	}
	app := container.NewContainer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.Bind(tt.args.abstract, tt.args.concrete, tt.args.shared)
			got := app.MakeWithParams(tt.args.abstract,tt.args.param)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWork() = %v, want %v", got, tt.want)
			}
		})
	}
}