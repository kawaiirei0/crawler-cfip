package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kawaiirei0/crawler-cfip/crawler"

	"github.com/gin-gonic/gin"
)

var (
	cache      []byte
	cacheMutex sync.RWMutex
	url        = "https://vps789.com/cfip/?remarks=ip"
)

func main() {
	// 启动定时抓取协程，每 1 小时更新一次
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			log.Println("开始抓取数据...")
			data, err := crawler.FetchTableData(url, 40*time.Second)
			if err != nil {
				log.Println("抓取失败:", err)
			} else {
				cacheMutex.Lock()
				cache = data
				cacheMutex.Unlock()
				log.Println("抓取完成，缓存更新")
			}
			<-ticker.C
		}
	}()

	// Gin 路由
	r := gin.Default()

	// API 接口返回缓存数据
	r.GET("/api/iplist", func(c *gin.Context) {
		cacheMutex.RLock()
		defer cacheMutex.RUnlock()
		if cache == nil {
			c.JSON(http.StatusOK, gin.H{"message": "数据尚未抓取"})
			return
		}
		c.Data(http.StatusOK, "application/json", cache)
	})

	log.Println("API 服务启动，监听 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
