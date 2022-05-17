package container

import (
	"fmt"
	"testing"
)

func BenchmarkAbstractToString(t *testing.B) {
	type args struct {
		abstract any
	}

	var itf *controllerInterface
	tests := []struct {
		name         string
		args         args
		wantResponse string
	}{
		{
			name:         "测试字符串",
			args:         args{abstract: "abc"},
			wantResponse: "abc",
		}, {
			name:         "测试数字类型",
			args:         args{abstract: 123},
			wantResponse: "abc",
		}, {
			name:         "测试方法",
			args:         args{abstract: addNum},
			wantResponse: "github.com/melodywen/go-ioc/container.addNum",
		}, {
			name:         "测试struct类型",
			args:         args{abstract: userController{}},
			wantResponse: "github.com/melodywen/go-ioc/container.userController",
		}, {
			name:         "测试指针类型",
			args:         args{abstract: newUserController()},
			wantResponse: "*github.com/melodywen/go-ioc/container.userController",
		}, {
			name:         "测试接口",
			args:         args{abstract: itf},
			wantResponse: "*github.com/melodywen/go-ioc/container.controllerInterface",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Println(err)
				}
			}()
			if gotResponse := AbstractToString(tt.args.abstract); gotResponse != tt.wantResponse {
				t.Errorf("AbstractToString() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func BenchmarkAbstractToString2(t *testing.B) {
	type args struct {
		identifier string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "如果是方法",
			args: args{identifier: AbstractToString(newUserController)},
			want: true,
		}, {
			name: "如果是数字",
			args: args{identifier: AbstractToString(addNum)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.B) {
			if got := isConstructMethod(tt.args.identifier); got != tt.want {
				t.Errorf("isConstructMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
