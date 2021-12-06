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
// Animal is a test struct
type Animal struct {
	name     string
	age      int
	category string
}

// NewAnimal animal construct
func NewAnimal(name string, age int, category string) *Animal {
	return &Animal{name: name, age: age, category: category}
}

// Bind a constructor
container.Bind(mock.Animal{}, mock.NewAnimal, false)

// resolve
container.MakeWithParams(mock.Animal{}, []interface{}{"dog", 12, "cate-pet"})

// Bind a callback method
container.Bind(mock.Animal{}, func() *mock.Animal {
        return mock.NewAnimal("dog", 12, "cate-pet")
    }, true)
// resolve
container.Make(mock.Animal{})

```
### 3. Binding of singletons
```golang
container.Singleton(mock.Animal{}, func() *mock.Animal {
    return mock.NewAnimal("dog", 12, "cate-pet")
})
```
### 4. Binding instance
```golang
container.Instance(mock.Animal{}, mock.NewAnimal("dog", 12, "cate-pet"))
```

##  License
© jiawenChen, 2021~time.Now

Released under the MIT License