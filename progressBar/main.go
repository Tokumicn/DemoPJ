package main

import (
	"fmt"
	"time"
)

func main() {
	var bar Bar
	bar.NewOption(0, 100)
	for i := 0; i <= 100; i++ {
		time.Sleep(100 * time.Millisecond)
		bar.Play(i)
	}

	bar.Finish()
}

type Bar struct {
	percent int    // 百分比
	cur     int    // 当前进度
	total   int    // 总进度
	rate    string // 进度条
	graph   string // 进度符号
}

func (bar *Bar) NewOption(start, total int) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}

	bar.percent = bar.getPercent()
	for i := 0; i < bar.percent; i += 2 {
		bar.rate += bar.graph // 初始化进度位置
	}
}

func (bar *Bar) getPercent() int {
	return int(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *Bar) NewOptionWithGraph(start, total int, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}

func (bar *Bar) Play(cur int) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()

	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%%  %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) Finish() {
	fmt.Println()
}
