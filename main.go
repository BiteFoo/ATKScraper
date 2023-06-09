package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BiteFoo/android_sdk_scraper/download"
	"github.com/BiteFoo/android_sdk_scraper/repos"
)

func main() {
	log.Println("Android Third-part Sdk Scraper v1.0.0")

	mgr := repos.NewLibManager()
	if mgr == nil {
		log.Panic("Initialize update error ")
	}
	log.Println("ctrl+c exit")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// start := time.Now()

	go download.PrintLog() // 格式化日志输出
	go mgr.Client.Run()
	mgr.DownloadATK()

	// elapse := time.Since(start)
	// log.Printf("共耗费时间: %v\n", elapse.Seconds())
	<-c
	log.Println("程序已经停止")

}
