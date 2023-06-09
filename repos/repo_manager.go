package repos

import (
	"log"
	"sync"

	"github.com/BiteFoo/android_sdk_scraper/download"
)

//主要入口

type LibManager struct {
	repoMap map[string]RepoLibrary
	Client  *download.DownloadClient
}

func NewLibManager() *LibManager {

	mgr := LibManager{
		Client:  download.NewDownloader(),
		repoMap: make(map[string]RepoLibrary, 0),
	}
	//全部下载
	mvnLibary := LoadLibraySpec("mvn-central-libraries.json")
	if mvnLibary != nil {
		mgr.repoMap["mvn-library"] = *mvnLibary
	}
	// 加载google
	googleLibraries := LoadLibraySpec("google-libraries.json")
	if googleLibraries != nil {
		log.Println("读取google-libraries完成")
		mgr.repoMap["google-library"] = *googleLibraries
	}
	amazonLibrary := LoadLibraySpec("amazon-libraries.json")
	if amazonLibrary != nil {
		log.Println("读取amazon-libraries完成")
		mgr.repoMap["amazon-library"] = *amazonLibrary
	}
	trackLibary := LoadLibraySpec("trackers.json")
	if trackLibary != nil {
		log.Println("读取trackers完成")
		mgr.repoMap["trackers-library"] = *trackLibary
	}
	//加载完成
	return &mgr
}

// 执行下载 第三方sdk并进行标记
func (mgr *LibManager) DownloadATK() {
	//直接执行下载
	// googleATK

	var wg sync.WaitGroup
	log.Println("start download atk")
	// wg.Add(1)
	for key, lib := range mgr.repoMap {

		switch key {
		case "mvn-library":
		case "amazon-library":
			// go DownloadMvnATK()
			wg.Add(1)
			log.Println("-> download ", key)
			go DownloadMvnATK(MVNBASE_URL, mgr.Client, lib)
		case "google-library":
			// case "trackers-library":
			wg.Add(1)
			go DownloadXMLLibs(mgr.Client, lib)
		}

		// }(key, lib)
	}
	wg.Wait()

}
