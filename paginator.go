package paginator

import (
	"bytes"
	"encoding/json"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type Pagination struct {
	Total        int         `json:"total"`
	Page         int         `json:"page"`
	Size         int         `json:"size"`
	Path         string      `json:"path"`
	PageCount    int         `json:"pageCount"`
	Slot         int         `json:"slot"`
	PrevName     string      `json:"-"`
	NextName     string      `json:"-"`
	PageList     []PageItem  `json:"-"`
	BarSize      string      `json:"-"` //默认是中等大小，pagination-sm是小，pagination-lg是大
	FirstPageUrl string      `json:"firstPageUrl"`
	LastPageUrl  string      `json:"lastPageUrl"`
	Queries      url.Values  `json:"-"`
	CustomTmpl   string      `json:"-"`
	Data         interface{} `json:"data"`
}
type PageItem struct {
	IsEllipsis bool // 是否为省略号
	PageNum    int  // 页码，仅在 IsEllipsis 为 false 时有效
}

func New(r *http.Request, total int, defalutSize ...int) *Pagination {
	var page, size int
	if len(defalutSize) > 0 {
		size = defalutSize[0]
	}
	values := r.URL.Query()
	if values.Has("page") {
		page, _ = strconv.Atoi(values.Get("page"))
	}
	p := Default(total, page)
	if values.Has("size") {
		size, _ = strconv.Atoi(values.Get("size"))
	}
	p.Path = r.URL.Path
	p.SetSize(size)
	return p

}

func Default(total int, page int) *Pagination {
	if total <= 0 {
		total = 0
	}
	if page <= 0 {
		page = 1
	}
	return &Pagination{
		Total:    total,
		Page:     page,
		Size:     10,
		Path:     "/",
		Slot:     5,
		PrevName: "«",
		NextName: "»",
	}
}

// func (p *Pagination) method() {}

func (p *Pagination) SetSize(size int) {
	if size <= 0 {
		size = 10
	}
	p.Size = size
}
func (p *Pagination) setPageCount() {
	p.PageCount = int(math.Ceil(float64(p.Total) / float64(p.Size)))
}

func (p *Pagination) HasPrevPage() bool {
	return p.Page != 1
}

func (p *Pagination) PrevLink() (link string) {

	if p.HasPrevPage() {
		link = p.GetLink(p.Page - 1)
	}
	return

}
func (p *Pagination) HasNextPage() bool {

	return p.Page != p.PageCount
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

func (p *Pagination) NextLink() (link string) {

	if p.HasNextPage() {
		page := p.Page + 1
		return p.GetLink(page)
	}
	return
}

func (p *Pagination) HasBar() bool {
	return p.PageCount > 0
}
func paginate(totalPages, currentPage, slot int) []PageItem {
	var pageList []PageItem

	// 确保 slot 为奇数，以便当前页居中
	if slot%2 == 0 {
		slot += 1
	}

	// 特殊情况：总页数少于等于 slot + 2 时，显示所有页码
	if totalPages <= slot+2 {
		for i := 1; i <= totalPages; i++ {
			pageList = append(pageList, PageItem{IsEllipsis: false, PageNum: i})
		}
		return pageList
	}

	// 计算分页槽的边界
	halfSlot := slot / 2
	start := currentPage - halfSlot
	end := currentPage + halfSlot

	// 确保边界在有效范围内
	if start < 2 { // 如果当前页接近首页
		start = 2
		end = start + slot - 1
	}
	if end > totalPages-1 { // 如果当前页接近尾页
		end = totalPages - 1
		start = end - slot + 1
	}

	// 添加首页
	pageList = append(pageList, PageItem{IsEllipsis: false, PageNum: 1})

	// 判断是否需要省略号（首页和槽的起点之间）
	if start > 2 {
		pageList = append(pageList, PageItem{IsEllipsis: true})
	}

	// 添加分页槽中的页码
	for i := start; i <= end; i++ {
		pageList = append(pageList, PageItem{IsEllipsis: false, PageNum: i})
	}

	// 判断是否需要省略号（槽的终点和尾页之间）
	if end < totalPages-1 {
		pageList = append(pageList, PageItem{IsEllipsis: true})
	}

	// 添加尾页
	pageList = append(pageList, PageItem{IsEllipsis: false, PageNum: totalPages})

	return pageList
}

func (p *Pagination) Paginate() {
	p.setPageCount()
	p.LastPageUrl = p.GetLink(p.PageCount)
	p.FirstPageUrl = p.GetLink(1)
	p.PageList = paginate(p.PageCount, p.Page, p.Slot)
}

func (p *Pagination) AddQuery(key string, value string) {
	if p.Queries == nil {
		p.Queries = make(url.Values)
	}
	if !p.Queries.Has(key) {
		p.Queries.Add(key, value)
	} else {
		p.Queries.Set(key, value)
	}
}

func (p *Pagination) AddQueries(values url.Values) {
	p.Queries = values
}

func (p *Pagination) GetLink(page int) (link string) {
	link = p.Path
	if page == 1 {
		if p.Queries.Has("page") {
			p.Queries.Del("page")
		}
	} else {
		p.AddQuery("page", strconv.Itoa(page))
	}
	if len(p.Queries) > 0 {
		link = p.Path + "?" + p.Queries.Encode()
	}
	return
}

// func (p *Pagination) Parse(tmpl string) string {
// 	var buf bytes.Buffer
// 	p.Paginate()
// 	t := template.Must(template.ParseFS(templates.Fs, "*.tmpl"))
// 	t.ExecuteTemplate(&buf, tmpl, map[string]interface{}{
// 		"p": p,
// 	})

// 	return buf.String()

// }

func (p *Pagination) ParseString(str string) string {
	var buf bytes.Buffer
	p.Paginate()
	t := template.Must(template.New("").Parse(str))
	t.Execute(&buf, map[string]interface{}{
		"p": p,
	})

	return buf.String()

}

// 1 是原生html
// 2 是bootstarp 简单模式，只有上一页当前页和下一页
// 3 是bootstrap 有分页列表模式
// 4 自定义模式, 通过p.CustomTmpl设置
func (p *Pagination) GetContent(t int) template.HTML {
	var tmpl string

	switch t {
	case 0:
		tmpl = p.CustomTmpl
	case 1:
		tmpl = `{{ if .p.HasBar}} <ul style="list-style: none;clear: both;"> {{if .p.HasPrevPage }} <li style="float: left;padding: 5px 10px; border: 1px solid gray; margin: 5px 5px;"> <a href="{{ .p.PrevLink }}">{{.p.PrevName}}</a> </li> {{else}} <li style="float: left;padding: 5px 10px; border: 1px solid gray; margin: 5px 5px;" aria-disabled="true"> <span>{{.p.PrevName}}</span> </li> {{end}} <li style="float: left;padding: 5px 10px; border: 1px solid gray; margin: 5px 5px;"> <a href="">{{.p.Page}}</a> </li> {{if .p.HasNextPage }} <li style="float: left;padding: 5px 10px; border: 1px solid gray; margin: 5px 5px;"> <a href="{{ .p.NextLink }}">{{.p.NextName}}</a> </li> {{else}} <li style="float: left;padding: 5px 10px; border: 1px solid gray; margin: 5px 5px;"> <span>{{.p.NextName}}</span> </li> {{end}} </ul> {{end}}`
	case 2:
		tmpl = `{{if .p.HasBar}} <nav aria-label="pagination"> <ul class="pagination"> {{if .p.HasPrevPage }} <li class="page-item"> <a href="{{ .p.PrevLink }}" class="page-link">{{.p.PrevName}}</a> </li> {{else}} <li class="page-item disabled"> <span class="page-link">{{.p.PrevName}}</span> </li> {{end}} <li class="page-item active" aria-current="page"> <span class="page-link">{{.p.Page}}</span> </li> {{if .p.HasNextPage }} <li class="page-item"> <a href="{{ .p.NextLink }}" class="page-link">{{.p.NextName}}</a> </li> {{else}} <li class="page-item disabled"> <span class="page-link">{{.p.NextName}}</span> </li> {{end}} </ul> </nav> {{end}}`
	case 3:
		tmpl = `{{ if .p.HasBar}} <nav aria-label="pagination"> <ul class="pagination {{.p.BarSize}}"> {{if .p.HasPrevPage }} <li class="page-item"> <a href="{{ .p.PrevLink }}" class="page-link">{{.p.PrevName}}</a> </li> {{else}} <li class="page-item disabled"> <span class="page-link">{{.p.PrevName}}</span> </li> {{end}} {{range $k,$v:= .p.PageList}} {{if $v.IsEllipsis}} <li class="page-item" aria-current="page"> <span class="page-link">...</span> </li> {{else}} {{if eq $.p.Page $v.PageNum}} <li class="page-item active" aria-current="page"><a class="page-link" href="{{$.p.GetLink $v.PageNum}}">{{$v.PageNum}}</a></li> {{else}} <li class="page-item"><a class="page-link" href="{{$.p.GetLink $v.PageNum}}">{{$v.PageNum}}</a></li> {{end}} {{end}} {{end}} {{if .p.HasNextPage }} <li class="page-item"> <a href="{{ .p.NextLink }}" class="page-link">{{.p.NextName}}</a> </li> {{else}} <li class="page-item disabled"> <span class="page-link">{{.p.NextName}}</span> </li> {{end}} </ul> </nav> {{end}}`
	case 4:
		tmpl = `{{ if .p.HasBar}} <nav aria-label="pagination"> <ul class="pagination {{.p.BarSize}}"> {{if .p.HasPrevPage }} <li class="page-item"> <a href="{{ .p.PrevLink }}" class="page-link">{{.p.PrevName}}</a> </li> {{else}} <li class="page-item mx-1 disabled"> <span class="page-link">{{.p.PrevName}}</span> </li> {{end}} {{range $k,$v:= .p.PageList}} {{if $v.IsEllipsis}} <li class="page-item mx-1" aria-current="page"> <span class="page-link">...</span> </li> {{else}} {{if eq $.p.Page $v.PageNum}} <li class="page-item active" aria-current="page"><a class="page-link" href="{{$.p.GetLink $v.PageNum}}">{{$v.PageNum}}</a></li> {{else}} <li class="page-item"><a class="page-link" href="{{$.p.GetLink $v.PageNum}}">{{$v.PageNum}}</a></li> {{end}} {{end}} {{end}} {{if .p.HasNextPage }} <li class="page-item"> <a href="{{ .p.NextLink }}" class="page-link">{{.p.NextName}}</a> </li> {{else}} <li class="page-item disabled"> <span class="page-link">{{.p.NextName}}</span> </li> {{end}} </ul> </nav> {{end}}`
	}
	return template.HTML(p.ParseString(tmpl))

}

func (p *Pagination) GetJson() string {
	p.Paginate()
	b, _ := json.Marshal(p)
	return string(b)
}
