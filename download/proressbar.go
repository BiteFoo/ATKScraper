package download

import (
	"fmt"
	"strings"
)

type Progressbar struct {
	title      string // 显示的进度条信息
	size       int64
	total      int64
	current    int64
	percentage int
	barLength  int
}

func NewBar(name string, barLength int, size int64) *Progressbar {
	return &Progressbar{
		title:     name,
		size:      size,
		total:     int64(barLength),
		barLength: barLength,
	}
}

func (p *Progressbar) Write(b []byte) (int, error) {

	n := len(b)
	p.current += int64(n)
	percentage := int(float32(p.current) / float32(p.size) * 100.0)

	if p.percentage < percentage {

		p.percentage = percentage
		completed := int(float32(p.current) / float32(p.size) * float32(p.barLength))
		remaing := p.barLength - completed
		if remaing <= 0 {
			//此时我们已经到达了相同的结果，因此这里直接填充满足100
			bar := "Downloading " + p.title + " [" + strings.Repeat("=", completed) + ">" + strings.Repeat(" ", 1) + "]"
			fmt.Printf("\r%s %d%%\n", bar, 100)
		} else {
			bar := "Downloading " + p.title + " [" + strings.Repeat("=", completed) + ">" + strings.Repeat(" ", remaing-1) + "]"
			fmt.Printf("\r%s %d%%", bar, p.percentage)
		}

	}
	return n, nil

}
