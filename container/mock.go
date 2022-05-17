package container

import "fmt"

func addNum(a int, b int) int {
	return a + b
}

func addAndParam(a int, b int) (int, int, int) {
	return a + b, a, b
}

type controllerInterface interface {
}

type userController struct {
}

func newUserController() *userController {
	return &userController{}
}
func newUserControllerAndOss(container *Container) (*userController, ossInterface) {
	var o *ossInterface
	return &userController{}, container.Make(o).(ossInterface)
}

type fileController struct {
}

func newFileController() *fileController {
	return &fileController{}
}

func newFileControllerAndOss(container *Container) (*fileController, ossInterface) {
	var o *ossInterface
	return &fileController{}, container.Make(o).(ossInterface)
}

func newUserControllerAndOther(request *Request, response Response, str string, i int) (*userController, Request, *Response, string, int) {
	return &userController{}, *request, &response, str, i
}

func newUserControllerAndObj(request *Request, response Response) (*userController, Request, *Response) {
	return &userController{}, *request, &response
}

type Request struct {
	Method string
	Uri    string
	param  map[string]any
}

func newRequest() *Request {
	return &Request{
		Method: "get",
		Uri:    "/user",
		param: map[string]any{
			"user": "张三",
			"age":  12,
		},
	}
}

type Response struct {
	Code  int
	Msg   string
	param map[string]any
}

func newResponse() *Response {
	return &Response{
		Code: 0,
		Msg:  "success",
		param: map[string]any{
			"user": "张三",
			"age":  12,
		},
	}
}

type mysql struct {
	Host string
	Port int
}

func newMysql() *mysql {
	return &mysql{Host: "localhost", Port: 3306}
}

type gorm struct {
	Adapter any
}

func newGorm(adapter mysql) *gorm {
	return &gorm{Adapter: adapter}
}

type ossInterface interface {
	Connect() string
}

type oss struct {
	Host string
	Port int
	Note string
}

func newOss() *oss {
	return &oss{
		Host: "localhost",
		Port: 7000,
		Note: "oss",
	}
}
func (o *oss) Connect() string {
	return fmt.Sprintf("host: %s,port: %d,note:%s", o.Host, o.Port, o.Note)
}

type ossAli struct {
	Host string
	Port int
	Note string
}

func newOssAli() *ossAli {
	return &ossAli{
		Host: "localhost",
		Port: 9000,
		Note: "ali oss",
	}
}
func (o *ossAli) Connect() string {
	return fmt.Sprintf("host: %s,port: %d,note:%s", o.Host, o.Port, o.Note)
}

type ossTencent struct {
	Host string
	Port int
	Note string
}

func newOssTencent() *ossTencent {
	return &ossTencent{
		Host: "localhost",
		Port: 8000,
		Note: "tencent oss",
	}
}
func (o *ossTencent) Connect() string {
	return fmt.Sprintf("host: %s,port: %d,note:%s", o.Host, o.Port, o.Note)
}
