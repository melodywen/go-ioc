package contracts

import (
	"fmt"
	container "github.com/melodywen/go-ioc"
	"testing"
)

func TestContainerContract(t *testing.T) {

	var app ContainerContract

	app = container.NewContainer()

	fmt.Println(app)
}
