package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	gq "github.com/PuerkitoBio/goquery"
)

//拉取依赖 go get github.com/PuerkitoBio/goquery

func fetch(code, name string) {
	if !strings.HasPrefix(code, "sz") && !strings.HasPrefix(code, "sh") {
		fmt.Printf("code要以sz或者sh开头: %s\n", code)
		return
	}

	url := fmt.Sprintf("https://qt.gtimg.cn/q=%s", code)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	doc, err := gq.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	data := strings.Split(strings.Split(strings.Trim(doc.Text(), "\";\n"), "=")[1], "~")
	if name == "" {
		name = data[1]
	}
	zs, current, incr := data[4], data[3], data[32]
	fmt.Printf("[%s(%s)] 昨收: %s, 当前: %s, 涨幅: %s%%\n", name, code, zs, current, incr)
}

func main() {
	for {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		fetch("sh600036", "招商银行")
		fetch("sz002722", "物产金轮")
		fetch("sh000001", "上证指数")
		time.Sleep(10 * time.Second) // 间隔10秒进行下一次请求
	}
}
