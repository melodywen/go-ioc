package container

import (
	"cjw.com/melodywen/go-ioc/mock"
	"fmt"
	"reflect"
	"testing"
)

func TestContainer_getConcrete(t *testing.T) {
	type fields struct {
		abstract interface{}
		concrete interface{}
		shared   bool
	}
	type args struct {
		abstract string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantConcrete interface{}
	}{
		// TODO: Add test cases.
		{
			name: "第一轮测试,测试对象",
			fields: fields{
				abstract: mock.Animal{},
				concrete: mock.NewAnimal("dog", 18, "cate"),
				shared:   false,
			},
			args:         args{abstract: "cjw.com/melodywen/go-ioc/mock.Animal"},
			wantConcrete: mock.NewAnimal("dog", 18, "cate"),
		},
		{
			name: "第一轮测试,测试指指针",
			fields: fields{
				abstract: &mock.Animal{},
				concrete: *mock.NewAnimal("dog", 18, "cate"),
				shared:   false,
			},
			args:         args{abstract: "*cjw.com/melodywen/go-ioc/mock.Animal"},
			wantConcrete: *mock.NewAnimal("dog", 18, "cate"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := newContainer()
			container.Bind(tt.fields.abstract, tt.fields.concrete, tt.fields.shared)
			if gotConcrete := container.getConcrete(tt.args.abstract); !reflect.DeepEqual(gotConcrete, tt.wantConcrete) {
				t.Errorf("getConcrete() = %v, want %v", gotConcrete, tt.wantConcrete)
			}
		})
	}
}

func TestContainer_IsShared(t *testing.T) {
	type fields struct {
		abstract interface{}
		concrete interface{}
		shared   bool
	}
	type args struct {
		abstract string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "测试——非分享",
			fields: fields{
				abstract: mock.NewAnimal,
				concrete: 1,
				shared:   false,
			},
			args: args{abstract: "cjw.com/melodywen/go-ioc/mock.NewAnimal"},
			want: false,
		},
		{
			name: "测试——分享",
			fields: fields{
				abstract: mock.NewAnimal,
				concrete: 1,
				shared:   true,
			},
			args: args{abstract: "cjw.com/melodywen/go-ioc/mock.NewAnimal"},
			want: true,
		}, {
			name: "测试——不存在",
			fields: fields{
				abstract: mock.NewAnimal,
				concrete: 1,
				shared:   true,
			},
			args: args{abstract: "abcc"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := newContainer()
			container.Bind(tt.fields.abstract, tt.fields.concrete, tt.fields.shared)
			if got := container.IsShared(tt.args.abstract); got != tt.want {
				fmt.Println(container)
				t.Errorf("IsShared() = %v, want %v", got, tt.want)
			}
		})
	}
}
