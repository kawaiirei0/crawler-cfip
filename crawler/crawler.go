package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"runtime"
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

// getChromePath 自动检测系统平台的 Chrome/Chromium 路径
func getChromePath() (string, error) {
	switch runtime.GOOS {
	case "linux":
		paths := []string{"/usr/bin/chromium", "/usr/bin/google-chrome", "/snap/bin/chromium"}
		for _, p := range paths {
			if _, err := exec.LookPath(p); err == nil {
				return p, nil
			}
		}
	case "windows":
		paths := []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		}
		for _, p := range paths {
			if _, err := exec.LookPath(p); err == nil {
				return p, nil
			}
		}
	case "darwin":
		p := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
		if _, err := exec.LookPath(p); err == nil {
			return p, nil
		}
	}
	return "", errors.New("Chrome/Chromium not found")
}

// FetchTableData 打开 URL 并抓取表格数据，返回 JSON bytes
func FetchTableData(url string, timeout time.Duration) ([]byte, error) {
	chromePath, err := getChromePath()
	if err != nil {
		return nil, err
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Headless, // 可选：添加无头模式
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
	)

	// 1️⃣ 创建执行器
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// 2️⃣ 创建浏览器上下文
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 3️⃣ 超时控制
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	var html string
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.el-table__row`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	doc.Find("script, style").Remove()

	var results []RowData
	doc.Find("tr.el-table__row").Each(func(i int, s *goquery.Selection) {
		cols := s.Find("td")
		if cols.Length() >= 6 {
			results = append(results, RowData{
				IP:       strings.TrimSpace(cols.Eq(0).Text()),
				Dx:       strings.TrimSpace(cols.Eq(1).Text()),
				Yd:       strings.TrimSpace(cols.Eq(2).Text()),
				Lt:       strings.TrimSpace(cols.Eq(3).Text()),
				Speed:    strings.TrimSpace(cols.Eq(4).Text()),
				LastTime: strings.TrimSpace(cols.Eq(5).Text()),
			})
		}
	})

	return json.MarshalIndent(results, "", "  ")
}
