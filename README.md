# go-ioc

![build status](https://camo.githubusercontent.com/f05145ad1c938e873697d2b624764921913522654e41fb7c68ba7918967a846b/68747470733a2f2f676f7265706f7274636172642e636f6d2f62616467652f6769746875622e636f6d2f676f2d676f726d2f676f726d)


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



## 单元测试结果：
生成覆盖率
```
go test -coverprofile=coverage.out

coverage: 92.2% of statements
ok      cjw.com/melodywen/go-ioc        0.219s
```
测试结果：
```
go tool cover -html=coverage.out 
```