package container

import (
	"github.com/melodywen/go-ioc/mock"
	"testing"
)

func TestCommonOfContainer_AbstractToString(t *testing.T) {
	type args struct {
		abstract interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "测试字符串",
			args: args{abstract: "abc"},
			want: "abc",
		}, {
			name: "测试int类型",
			args: args{abstract: 1234},
			want: "1234",
		}, {
			name: "测试struct类型",
			args: args{abstract: mock.Animal{}},
			want: "github.com/melodywen/go-ioc/mock.Animal",
		}, {
			name: "测试指针类型",
			args: args{abstract: mock.NewAnimal("猫", 1, "猫科")},
			want: "*github.com/melodywen/go-ioc/mock.Animal",
		}, {
			name: "测试指针类型",
			args: args{abstract: mock.NewAnimal},
			want: "github.com/melodywen/go-ioc/mock.NewAnimal",
		}, {
			name: "测试指针类型-对象方法",
			args: args{abstract: (&mock.Animal{}).Say},
			want: "github.com/melodywen/go-ioc/mock.(*Animal).Say-fm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extend := &ExtendOfContainer{}
			if got := extend.AbstractToString(tt.args.abstract); got != tt.want {
				t.Errorf("AbstractToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
