package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

func main() {
	myhash := New(3, nil)

	myhash.Add("A")
	myhash.Add("B")
	myhash.Add("C")

	x := 0
	for x < 50 {
		fmt.Println(myhash.Get("111132231231"))
		x++
	}

}

type Hash func(data []byte) uint32
type Map struct {
	hash     Hash           // 计算Hash值的函数
	replicas int            // 副本数，这里影响到虚拟节点的个数
	keys     []int          // 有序列表，从大到小排序
	hashMap  map[int]string // 记录虚拟节点和物理节点元数据关系
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		// 默认使用crc32计算hash
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// keys ==> [A, B, C]
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// hash值 = hash(i+key)
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)

			// 虚拟节点 <--> 实际节点
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}

	// 计算用户输入key的hash值
	hash := int(m.hash([]byte(key)))
	// 查看落在那个值域范围？选择虚拟节点
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	if idx == len(m.keys) {
		idx = 0 // 环
	}

	index := m.keys[idx]
	// 选择对应屋里节点
	return m.hashMap[index]
}

func (m *Map) IsEmpty() bool {
	return false
}
