package contracts

import (
	container "cjw.com/melodywen/go-ioc"
	"fmt"
	"testing"
)


func TestContainerContract(t *testing.T) {
	
	var app ContainerContract;
	
	app = container.NewContainer()
	
	fmt.Println(app)
}
