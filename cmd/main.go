package main

import (
	"fmt"

	"github.com/wujiyu98/paginator"
)

// func getPages(currentPage uint, total uint, slot uint, perPage uint) {
// 	var lists []string
// 	pageCount := uint(math.Ceil(float64(total) / float64(perPage)))
// 	switch {
// 	case pageCount <= slot:
// 		for i := 1; i <= int(pageCount); i++ {
// 			lists = append(lists, fmt.Sprint(i))
// 		}
// 	default:
// 		switch {
// 		case currentPage > pageCount-(slot+1)/2:
// 			lists = append(lists, `1`)
// 			lists = append(lists, `...`)
// 			for i := int(pageCount - (slot - 2)); i <= int(pageCount); i++ {
// 				lists = append(lists, fmt.Sprint(i))
// 			}

// 		case currentPage <= (slot+1)/2:
// 			for i := 1; i <= int(slot-2); i++ {
// 				lists = append(lists, fmt.Sprint(i))
// 			}
// 			lists = append(lists, `...`)
// 			lists = append(lists, fmt.Sprint(pageCount))
// 			fmt.Println("c2")
// 		default:
// 			lists = append(lists, `1`)
// 			lists = append(lists, `...`)
// 			for i := int(currentPage - (slot-2)/2); i <= int(currentPage+(slot-3)/2); i++ {
// 				lists = append(lists, fmt.Sprint(i))
// 			}
// 			lists = append(lists, `...`)
// 			lists = append(lists, fmt.Sprint(pageCount))
// 			fmt.Println("c3")

// 		}

// 	}

// 	fmt.Println(lists)

// }

func main() {

	p := paginator.Default(100, 10)
	p.CurrentPage = 7

	p.AddQuery("keyword", "hello")
	p.AddQuery("age", "31231")

	p.GetLists()
	fmt.Println(p.GetLastPageUrl())

}
