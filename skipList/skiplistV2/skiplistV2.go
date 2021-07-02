package skiplistV2

import (
	"fmt"
	"math/rand"
)

func main() {
	sl := New()

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

	fmt.Println(sl.len)
	e = sl.Delete(float64(50))
	fmt.Println(e.Value)
	fmt.Println(sl.len)

	for e := sl.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

}

const (
	maxLevel int     = 16
	p        float32 = 0.25
)

type Element struct {
	Score   float64
	Value   interface{}
	forward []*Element
}

func newElement(score float64, value interface{}, level int) *Element {
	return &Element{
		Score:   score,
		Value:   value,
		forward: make([]*Element, level),
	}
}

type SkipListV2 struct {
	header *Element
	len    int
	level  int
}

func New() *SkipListV2 {
	return &SkipListV2{
		header: &Element{
			forward: make([]*Element, maxLevel),
		},
	}
}

func (sl *SkipListV2) GetLength() int {
	return sl.len
}

func randomLevel() int {
	level := 1
	for rand.Float32() < p && level < maxLevel {
		level++
	}
	return level
}

func (sl *SkipListV2) Front() *Element {
	return sl.header.forward[0]
}

func (e *Element) Next() *Element {
	if e != nil {
		return e.forward[0]
	}
	return nil
}

func (sl *SkipListV2) Search(score float64) (element *Element, ok bool) {
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].Score < score {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	if x != nil && x.Score == score {
		return x, true
	}

	return nil, false
}

func (sl *SkipListV2) Insert(score float64, value interface{}) *Element {
	update := make([]*Element, maxLevel)
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].Score < score {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]

	if x != nil && x.Score == score {
		x.Value = value
		return x
	}

	level := randomLevel()
	// 新增一层
	if level > sl.level {
		level = sl.level + 1
		update[sl.level] = sl.header
		sl.level = level
	}

	e := newElement(score, value, level)
	for i := 0; i < level; i++ {
		e.forward[i] = update[i].forward[i]
		update[i].forward[i] = e
	}
	sl.len++
	return e
}

func (sl *SkipListV2) Delete(score float64) *Element {
	update := make([]*Element, maxLevel)
	x := sl.header

	// 找到节点
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].Score < score {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]

	// 删除节点
	if x != nil && x.Score == score {
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] != x {
				return nil
			}
			update[i].forward[i] = x.forward[i]
		}
		sl.len--
	}
	return x
}
