package main

import (
	"log"
	"time"

	"github.com/BiteFoo/android_sdk_scraper/download"
	"github.com/BiteFoo/android_sdk_scraper/repos"
)

func main() {
	log.Println("Android Third-part Sdk Scraper v1.0.0")

	mvn := repos.NewMvn()
	if mvn == nil {
		log.Panic("Initialize update error ")
	}
	log.Println("ctrl+c exit")

	start := time.Now()
	go download.PrintLog() // 格式化日志输出
	go repos.GetReposInfo(mvn)
	mvn.Client.Run()
	mvn.Client.Wait()
	elapse := time.Since(start)
	log.Printf("共耗费时间: %v\n", elapse.Seconds())

}
