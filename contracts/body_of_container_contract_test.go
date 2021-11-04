package contracts

import (
	container "cjw.com/melodywen/go-ioc"
	"fmt"
	"testing"
)

func Test_BodyOfContainerContract(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "测试是否实现接口",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var build BodyOfContainerContract
			build = new(container.BodyOfContainer)
			fmt.Println(build)
		})
	}
}
