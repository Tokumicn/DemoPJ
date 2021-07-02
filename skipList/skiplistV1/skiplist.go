package skiplistV1

import (
	"fmt"
	"math/rand"
	"time"
)

// [已理解] 后续将特点整理成文章
func main() {

	obj := Construtor()
	search := obj.Search(111)
	fmt.Println(search)

	obj.Add(111)
	search2 := obj.Search(111)
	fmt.Println(search2)

	erase := obj.Erase(111)
	fmt.Println(erase)
}

const MAX_LEVEL = 10

type SkipList struct {
	maxLevel int
	root     []*Node
}

type Node struct {
	level int
	next  []*Node
	prev  []*Node
	value int
	count int
}

func NewNode(level, value int) *Node {
	return &Node{
		level: level,
		value: value,
		next:  make([]*Node, level+1),
		prev:  make([]*Node, level+1),
		count: 1,
	}
}

func Construtor() SkipList {
	return SkipList{
		maxLevel: -1,
	}
}

// 查找并返回节点
func (this *SkipList) SearchPos(target int) *Node {

	if this.maxLevel == -1 {
		return nil
	}

	// 从最上层开始搜索,最上层元素最少
	next := this.root[this.maxLevel]
	for level := next.level; level >= 0; level-- {
		if target == next.value {
			return next
		} else if target < next.value {
			for next.prev[level] != nil && target <= next.prev[level].value {
				next = next.prev[level]
			}
		} else {
			for next.next[level] != nil && target >= next.next[level].value {
				next = next.next[level]
			}
		}
		fmt.Println("search :", target, " level :", level)
	}
	return next
}

// 查找是否存在
func (this *SkipList) Search(target int) bool {
	node := this.SearchPos(target)

	if node == nil {
		return false
	}
	return node.value == target
}

func (this *SkipList) Print() {
	if this.root == nil {
		return
	}

	fmt.Println(*this)

	for i := this.maxLevel; i >= 0; i-- {
		cur := this.root[i]
		for cur != nil && cur.next != nil {
			fmt.Print("->[", cur.value, cur.count, "]")
			cur = cur.next[i]
		}
		fmt.Println(i)
	}
}

func (this *SkipList) randomUpgrade() bool {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(7) % 2

	if this.maxLevel > MAX_LEVEL {
		return false
	}
	fmt.Println("rand: ----", r)
	return r == 1
}

func (n *Node) InsertLevelNode(level int, newNode *Node) {
	if n == nil || n.prev == nil || n.next == nil || newNode == nil || newNode.prev == nil || newNode.next == nil {
		return
	}

	// 插入值小于当前节点，因此插入点在前方
	if n.value > newNode.value {
		prev := n.prev[level]

		n.prev[level] = newNode
		newNode.next[level] = n
		newNode.prev[level] = prev

		if prev != nil {
			prev.next[level] = newNode
		}
	} else {
		next := n.next[level]

		n.next[level] = newNode
		newNode.prev[level] = n
		newNode.next[level] = next

		if next != nil {
			next.prev[level] = newNode
		}
	}
}

// 如果元素x出现在i层，则所有小于i的层都包含x
func (this *SkipList) Add(num int) {
	n := this.SearchPos(num)

	// 元素不存在，加入
	if n == nil {
		this.maxLevel = 0
		n = NewNode(0, num)
		this.root = append(this.root, n)
		return
	}

	// 元素存在，则计数器加1
	if n.value == num {
		n.count++
		return
	}

	if this.randomUpgrade() {
		this.maxLevel++
		nn := NewNode(this.maxLevel, num)
		this.root = append(this.root, nn)

		for i := 1; i <= this.maxLevel-1; i++ {
			in := this.root[i]
			for in != nil && in.value > num && in.prev != nil && in.prev[i] != nil {
				in = in.prev[i]
			}

			for in != nil && in.value < num && in.next != nil && in.next[i] != nil {
				in = in.next[i]
			}

			// 从第一次开始，之后的每一层都要有这个元素
			in.InsertLevelNode(i, nn)
		}
		// 第一层(最低一层)要包含全部元素
		n.InsertLevelNode(0, nn)
	} else {
		nn := NewNode(0, num)
		n.InsertLevelNode(0, nn)
	}
}

func (this *SkipList) DeleteNode(n *Node, level int) {
	if n == nil {
		return
	}

	next := n.next[level]
	prev := n.prev[level]

	if prev != nil {
		prev.next[level] = next
	} else {
		this.root[level] = next
	}

	if next != nil {
		next.prev[level] = prev
	}
}

func (this *SkipList) Erase(num int) bool {
	n := this.SearchPos(num)
	if n == nil || n.value != num {
		return false
	}

	if n.count > 1 {
		n.count--
		return true
	}

	// 删除时，每一层的该元素都要删除，当count>1时仅减少count,并不删除节点
	for i := 0; i <= n.level; i++ {
		this.DeleteNode(n, i)
	}

	count := 0
	for level := this.maxLevel; n.level == this.maxLevel && level >= 0 && this.root != nil && this.root[level] == nil; level-- {
		count++
	}

	this.root = this.root[:len(this.root)-count]
	this.maxLevel -= count

	n.level = -1
	n.next = nil
	n.prev = nil

	return true
}
