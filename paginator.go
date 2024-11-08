package paginator

import (
	"bytes"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"text/template"
)

type Pagination struct {
	Total     int         `json:"total"`
	Page      int         `json:"page"`
	Size      int         `json:"size"`
	Path      string      `json:"path"`
	PageCount int         `json:"pageCount"`
	Slot      int         `json:"slot"`
	PrevName  string      `json:"-"`
	NextName  string      `json:"-"`
	Queries   url.Values  `json:"-"`
	Data      interface{} `json:"data"`
}
type PageItem struct {
	IsEllipsis bool // 是否为省略号
	PageNum    int  // 页码，仅在 IsEllipsis 为 false 时有效
}

func New(total int, page int) *Pagination {
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
		PrevName: "&laquo;",
		NextName: "&raquo;",
	}
}

// func (p *Pagination) method() {}

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

func (p *Pagination) NextLink() (link string) {

	if p.HasNextPage() {
		page := p.Page + 1
		fmt.Println("-------", page)
		return p.GetLink(page)
	}
	return
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

func (p *Pagination) Paginate() (pageList []PageItem) {
	p.setPageCount()
	pageList = paginate(p.PageCount, p.Page, p.Slot)
	return
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

func (p *Pagination) GetLink(page int) (link string) {
	if page != 1 {
		p.AddQuery("page", strconv.Itoa(page))
		p.AddQuery("size", strconv.Itoa(p.Size))
	}
	link = p.Path
	if len(p.Queries) > 0 {
		link = p.Path + "?" + p.Queries.Encode()
	}

	return
}

func (p *Pagination) Parse() {
	p.Paginate()
	var buf bytes.Buffer
	t := template.Must(template.ParseGlob("template/*.tmpl"))
	err := t.ExecuteTemplate(&buf, "bs5.tmpl", map[string]interface{}{
		"p": p,
	})
	fmt.Println(err)
	fmt.Println(buf.String())

}
