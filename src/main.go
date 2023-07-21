package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetch(code, name string) (string, error) {
	if !strings.HasPrefix(code, "sz") && !strings.HasPrefix(code, "sh") {
		return "", fmt.Errorf("code要以sz或者sh开头: %s", code)
	}

	url := fmt.Sprintf("https://qt.gtimg.cn/q=%s", code)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error getting data: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error parsing HTML: %v", err)
	}

	data := strings.Split(strings.Split(strings.Trim(doc.Text(), "\";\n"), "=")[1], "~")
	if name == "" {
		name = data[1]
	}
	current := data[3]

	return fmt.Sprintf("value{stock_id=\"%s\",name=\"%s\"} %s", code, name, current), nil
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	result := ""
	stocks := [][]string{
		{"sh600036", "招商银行"},
		{"sz002722", "物产金轮"},
		{"sh000001", "上证指数"},
	}

	for _, stock := range stocks {
		metric, err := fetch(stock[0], stock[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result += metric + "\n"
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprint(w, result)
}

func main() {
	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(":8080", nil)
}

