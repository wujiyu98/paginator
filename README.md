# paginator 使用说明

### 引用
go get github.com/wujiyu98/paginator


r:= &http.Request{}

也可以通过gin r :=ctx.Request
默认size是10
p := paginator.New(r,1000, 20)
p.GetContent(3)

1. 选项0 是自定义模式, 通过p.CustomTmpl设置
2. 选项1 是HTML原生列表样式
3. 选项2 是bootstarp 简单模式，只有上一页当前页和下一页
4. 选项3 是bootstrap 有分页列表模式
5. 选项4 是在3的基础上添加间隔