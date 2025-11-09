package crawler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type RowData struct {
	IP       string `json:"ip"`
	Dx       string `json:"dx"`
	Yd       string `json:"yd"`
	Lt       string `json:"lt"`
	Speed    string `json:"speed"`
	LastTime string `json:"last_time"`
}

// FetchTableData 打开 URL 并抓取表格数据，返回 JSON bytes
func FetchTableData(url string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 设置抓取超时
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.el-table__row`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return nil, err
	}

	// 解析 HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// 删除无关标签
	doc.Find("script, style").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	var results []RowData
	doc.Find("tr.el-table__row").Each(func(i int, s *goquery.Selection) {
		cols := s.Find("td")
		if cols.Length() >= 6 {
			row := RowData{
				IP:       strings.TrimSpace(cols.Eq(0).Text()),
				Dx:       strings.TrimSpace(cols.Eq(1).Text()),
				Yd:       strings.TrimSpace(cols.Eq(2).Text()),
				Lt:       strings.TrimSpace(cols.Eq(3).Text()),
				Speed:    strings.TrimSpace(cols.Eq(4).Text()),
				LastTime: strings.TrimSpace(cols.Eq(5).Text()),
			}
			results = append(results, row)
		}
	})

	return json.MarshalIndent(results, "", "  ")
}
