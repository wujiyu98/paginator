package main

import (
	"fmt"

	"github.com/wujiyu98/paginator/v2"
)

// PageItem 结构体表示一个页码项，可以是数字页码，也可以是省略号
type PageItem struct {
	IsEllipsis bool // 是否为省略号
	PageNum    int  // 页码，仅在 IsEllipsis 为 false 时有效
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

func testPaginate() {
	p := paginator.New(100, 4)
	pageList := p.Paginate()

	// 打印分页结果
	for _, item := range pageList {
		if item.IsEllipsis {
			fmt.Print("... ")
		} else {
			fmt.Printf("%d ", item.PageNum)
		}
	}
	// 输出示例：1 2 3 4 ... 10
	p.Parse()

}

func main() {
	testPaginate()

}
