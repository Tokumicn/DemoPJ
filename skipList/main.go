package main

import (
	"DemoPJ/skipList/skiplistV2"
	"fmt"
)

func main() {
	sl := skiplistV2.New()

	sl.Insert(float64(100), "foo")

	e, ok := sl.Search(float64(100))
	fmt.Println(ok)
	fmt.Println(e.Value)

	e, ok = sl.Search(float64(200))
	fmt.Println(ok)
	fmt.Println(e)

	sl.Insert(float64(20.5), "bar")
	sl.Insert(float64(50), "spam")
	sl.Insert(float64(20), 42)

	fmt.Println(sl.GetLength())
	e = sl.Delete(float64(50))
	fmt.Println(e.Value)
	fmt.Println(sl.GetLength())

	for e := sl.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
