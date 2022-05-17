# go-ioc
[![Build Status](https://github.com/gin-gonic/gin/workflows/Run%20Tests/badge.svg?branch=master)]()
[![Go Report Card](https://camo.githubusercontent.com/f05145ad1c938e873697d2b624764921913522654e41fb7c68ba7918967a846b/68747470733a2f2f676f7265706f7274636172642e636f6d2f62616467652f6769746875622e636f6d2f676f2d676f726d2f676f726d)]()
[![GoDoc](https://pkg.go.dev/badge/github.com/gin-gonic/gin?status.svg)]()
[![license](https://camo.githubusercontent.com/992daabc2aa4463339825f8333233ba330dd08c57068f6faf4bb598ab5a3df2e/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f6c6963656e73652d4d49542d627269676874677265656e2e737667)]()

## Overview
Create an IOC component modeled after Laravel to address dependency injection and inversion of control issues. The main part of the basic comprehensive coverage

- The binding
- The singleton
- The instance
- The contextual binding map
- The registered type aliases.
- Be building to resolve parameter by contexture
- All of the global [before resolving/resolving/after resolving] callbacks
- All of the registered rebound callbacks
- "Extend" an abstract type in the container.

## Getting Started
### 1. get container
```golang
container := NewContainer()
```
### 2. Simple binding
```golang
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

container.Bind(&Request{}, newRequest, false)

container.Make(&Request{})

container.Bind(&Request{}, func() *Request {
    return newRequest()
}, true)
```
### 3. Binding of singletons
```golang
container.Singleton(&Request{}, func() *Request {
    return newRequest()
})
```
### 4. Binding instance
```golang
container.Instance(&Request{}, newRequest())
```
### 5 Bind interfaces to implementations
```golang
var cacheInterface *mock.cacheInterface
container.Singleton(cacheInterface, mock.NewRedisCache)
```
### 6 Context binding
> Refer to context binding in the Laravel documentation for details on context binding

```golang
var o *ossInterface
container.When([]any{newUserControllerAndOss}).Need(o).Give(func(oss *ossAli) ossInterface {
    return oss
})
container.When([]any{newFileControllerAndOss}).Need(o).Give(func(oss *ossTencent) ossInterface {
    return oss
})
```
### 7. Automatic injection
```golang
container.SingletonIf(&Request{}, newRequest)
container.Instance(Request{}, *newRequest())
container.SingletonIf(&Response{}, newResponse)
container.Instance(Response{}, *newResponse())


func newUserControllerAndObj(request *Request, response Response) (*userController, Request, *Response) {
	return &userController{}, *request, &response
}

container.Bind(newUserControllerAndObj, nil)
container.Make(newUserControllerAndObj)
```
### 8. Container events
```golang
container.BeforeResolving(nil, func(s string, param []any, c *Container) {
		fmt.Println("global-before-callback", s)
})
container.BeforeResolving(newUserControllerAndObj, func(s string, param []any, c *Container) {
    fmt.Println("before-callback", s)
})
container.Resolving(nil, func(instance any, container *Container) {
    fmt.Println("global-ing-callback", instance)
})
container.Resolving(newUserControllerAndObj, func(instance any, container *Container) {
    fmt.Println("ing-callback", instance)
})
container.AfterResolving(nil, func(instance any, container *Container) {
    fmt.Println("global-after-callback", instance)
})
container.AfterResolving(newUserControllerAndObj, func(instance any, container *Container) {
    fmt.Println("after-callback", instance)
})
container.Extend(newUserControllerAndObj, func(object any, container *Container) any {
    return object
})
container.Rebinding(newUserControllerAndObj, func(container *Container, instance any) {
    fmt.Println("rebind callback", AbstractToString(newUserControllerAndObj))
})
....
```
## Notice
The unit test coverage of this project is 100%, and most of the methods have specific test instances, except for the part of exception throwing errors. For details on each direction of learning, refer to the test code.

##  License
© jiawenChen, 2022~time.Now

Released under the MIT License