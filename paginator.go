package paginator

import (
	"fmt"
	"math"
	"strings"
)

type Paginate struct {
	Total       uint
	PageCount   uint
	PerPage     uint
	CurrentPage uint
	LastPage    uint
	Query       map[string]string
	Path        string
	Sort        string
	From        string
	Slot        uint
}

func Default(total uint, perpage uint) *Paginate {
	p := Paginate{
		Total:   total,
		PerPage: perpage,
	}
	p.Query = make(map[string]string)
	p.CurrentPage = 1
	p.Slot = 5
	p.Sort = "id-desc"
	p.Path = "/"
	p.setPageCount()
	return &p
}

func (p *Paginate) setPageCount() {
	p.PageCount = uint(math.Ceil(float64(p.Total) / float64(p.PerPage)))
	p.LastPage = p.PageCount
}

func (p *Paginate) valiCurrentPage() {
	if p.CurrentPage <= 0 {
		p.CurrentPage = 1
	}
	if p.CurrentPage > p.PageCount {
		p.CurrentPage = p.PageCount
	}
}

func (p *Paginate) GetLists() {
	var lists []string
	p.setPageCount()
	p.valiCurrentPage()
	switch {
	case p.PageCount <= p.Slot:
		for i := 1; i <= int(p.PageCount); i++ {
			lists = append(lists, fmt.Sprint(i))
		}
	default:
		switch {
		case p.CurrentPage > p.PageCount-(p.Slot+1)/2:
			lists = append(lists, `1`)
			lists = append(lists, `...`)
			for i := int(p.PageCount - (p.Slot - 2)); i <= int(p.PageCount); i++ {
				lists = append(lists, fmt.Sprint(i))
			}

		case p.CurrentPage <= (p.Slot+1)/2:
			for i := 1; i <= int(p.Slot-2); i++ {
				lists = append(lists, fmt.Sprint(i))
			}
			lists = append(lists, `...`)
			lists = append(lists, fmt.Sprint(p.PageCount))

		default:
			lists = append(lists, `1`)
			lists = append(lists, `...`)
			for i := int(p.CurrentPage - (p.Slot-2)/2); i <= int(p.CurrentPage+(p.Slot-3)/2); i++ {
				lists = append(lists, fmt.Sprint(i))
			}
			lists = append(lists, `...`)
			lists = append(lists, fmt.Sprint(p.PageCount))

		}
	}
	fmt.Println(lists)

}

func (p *Paginate) AddQuery(key string, value string) {

	p.Query[key] = value

}

func (p *Paginate) AddQueries(items map[string]string) {

	for k, v := range items {
		p.Query[k] = v

	}

}

func (p *Paginate) addUrlParam(path string, k string, v string) string {
	if !strings.Contains(path, "?") {
		path += fmt.Sprint("?", k, "=", v)
	} else {
		path += fmt.Sprint("&", k, "=", v)
	}
	return path
}

func (p *Paginate) GetQueryUrl() string {
	path := p.Path
	for k, v := range p.Query {
		path = p.addUrlParam(path, k, v)
	}
	return path
}

func (p *Paginate) GetFirstPageUrl() string {
	path := p.GetQueryUrl()
	return path
}

func (p *Paginate) setComonParam(path string) string {
	path = p.addUrlParam(path, "total", fmt.Sprint(p.Total))
	path = p.addUrlParam(path, "per-page", fmt.Sprint(p.PerPage))
	path = p.addUrlParam(path, "sort", p.Sort)
	return path

}

func (p *Paginate) GetLastPageUrl() string {
	path := p.GetQueryUrl()
	path = p.addUrlParam(path, "page", fmt.Sprint(p.PageCount))
	return p.setComonParam(path)
}
