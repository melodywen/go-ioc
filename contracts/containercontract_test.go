package contracts

import (
	container "cjw.com/melodywen/go-ioc"
	"fmt"
	"testing"
)

func Test_containerContract(t *testing.T) {

	var app ContainerContracts

	app = new(container.Container)

	fmt.Println(app)

}
