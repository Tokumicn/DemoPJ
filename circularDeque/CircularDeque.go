package main

// [未解决]  需要结合图形理解
func main() {

}

/// |head ======================== rear|

type MyCircularDeque struct {
	length int
	data   []int
	head   int
	rear   int
}

func Constructor(k int) MyCircularDeque {
	return MyCircularDeque{
		length: k + 1,
		data:   make([]int, k+1), // 空一格位置区分满和空
		head:   0,
		rear:   0,
	}
}

// 判断空
func (this *MyCircularDeque) IsEmpty() bool {
	return this.head == this.rear
}

// 判断满
func (this *MyCircularDeque) IsFull() bool {
	return (this.rear+1)%this.length == this.head
}

// 从头部插入
func (this *MyCircularDeque) InsertFront(value int) bool {
	if this.IsFull() {
		return false
	}

	if this.IsEmpty() {
		if this.rear == this.length-1 {
			// rear已经到了头部，循环到尾部
			this.rear = 0
		} else {
			// ？？？
			this.rear++
		}
		this.data[this.head] = value
		return true
	}

	if this.head == 0 {
		this.head = this.length - 1
	} else {
		this.head--
	}

	this.data[this.head] = value
	return true
}

func (this *MyCircularDeque) InsertLast(value int) bool {
	if this.IsFull() {
		return false
	}

	//if this.IsEmpty() {
	//	this.data[this.rear] = value
	//	if this.rear == this.length-1 {
	//		// 最后一个空间使用后，循环到头
	//		this.rear = 0
	//	} else {
	//		this.rear++
	//	}
	//	return true
	//}

	this.data[this.rear] = value
	if this.rear == this.length-1 {
		// 最后一个空间使用后，循环到头
		this.rear = 0
	} else {
		this.rear++
	}
	return true
}

func (this *MyCircularDeque) DelelteFront() bool {
	if this.IsEmpty() {
		return false
	}

	if this.head == this.length-1 {
		this.head = 0
	} else {
		this.head++
	}
	return true
}

func (this *MyCircularDeque) DelelteLast() bool {
	if this.IsEmpty() {
		return false
	}

	if this.rear == 0 {
		this.rear = this.length - 1
	} else {
		this.rear--
	}
	return true
}

func (this *MyCircularDeque) GetFront() int {
	if this.IsEmpty() {
		return -1
	}
	return this.data[this.head]
}

func (this *MyCircularDeque) GetRear() int {
	if this.IsEmpty() {
		return -1
	}

	if this.rear == 0 {
		return this.data[this.length-1]
	}
	return this.data[this.rear-1]
}
