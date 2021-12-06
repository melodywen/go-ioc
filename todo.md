## 常用命令
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

常用工具：
```
https://github.com/gojp/goreportcard
```