package util

import (
	"bytes"
	"os"
	"testing"
)

func TestParseUrlPage(t *testing.T) {
	jekyllDir := `D:\code\dejavuzhou.github.io`
	err := ParseUrlPage("https://segmentfault.com/a/1190000015591319", "div.article__content", jekyllDir)
	if err != nil {
		t.Error(err)
	}
}

func TestConvert(t *testing.T) {
	h := `
<ol>
<li>减少<code>gc</code>压力，栈上的变量，随着函数退出后系统直接回收，不需要<code>gc</code>标记后再清除。</li>
<li>减少内存碎片的产生。</li>
<li>减轻分配堆内存的开销，提高程序的运行速度。</li>
</ol>
<h2 id="articleHeader4">如何确定是否逃逸</h2>
`
	var inBuff = bytes.NewBuffer([]byte(h))

	err := convert(os.Stdout, inBuff, &Option{})
	if err != nil {
		t.Error(err)
	}
}
