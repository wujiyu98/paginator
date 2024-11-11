# paginator
go paginator

go get github.com/wujiyu98/paginator/v2

使用 
r:= &http.Request{}
//也可以通过gin r :=ctx.Request
//默认size是10
p := paginator.New(r,1000)
p.GetContent(3)

// 1 是原生html
// 2 是bootstarp 简单模式，只有上一页当前页和下一页
// 3 是bootstrap 有分页列表模式
// 4 自定义模式, 通过p.CustomTmpl设置