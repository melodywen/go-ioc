# go-ioc
自创一个go 的ioc

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