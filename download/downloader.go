package download

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BiteFoo/android_sdk_scraper/utils"
)

//下载器接口

type DownloadClient struct {
	Wg            sync.WaitGroup
	DownloadChans chan DownloadInfo
}

// 下载信息
type DownloadInfo struct {
	Category    string //  标识sdk类别
	LibName     string //  libName
	LibVersion  string //  lib version
	FileName    string //  标识sdk名的例如 audience-network-sdk-6.14.0.aar
	DownloadURL string //  下载url
	SaveRoot    string // 保存路径
	ReleaseDate string // 发布时间
	Comment     string
}

type LibXmlTemplate struct {
	XMLName        xml.Name `xml:"library"`
	LibName        string   `xml:"name"`
	LibCategory    string   `xml:"category"`
	LibVersion     string   `xml:"version"`
	LibReleaseDate string   `xml:"release"`
	LibComment     string   `xml:"comment"`
}

const (
	MAX_POOL = 10 //同时支持10个协程执行即可
)

var (
	logChan = make(chan string)
) // 保证日志输出结果

func NewDownloader() *DownloadClient {
	return &DownloadClient{
		Wg:            sync.WaitGroup{},
		DownloadChans: make(chan DownloadInfo),
	}
}

func (client *DownloadClient) Run() {
	//

	for i := 0; i < MAX_POOL; i++ {
		client.Wg.Add(1)
		go func() {
			defer client.Wg.Done()
			for download := range client.DownloadChans {
				//client.Wg.Add(1)
				go saveFile(download)
			}
		}()
	}
}

func (c *DownloadClient) Submit(info DownloadInfo) {
	c.DownloadChans <- info
}

func (c *DownloadClient) Wait() {
	c.Wg.Wait() //等待结束
}

func saveFile(info DownloadInfo) {
	// defer wg.Done()

	save := info.SaveRoot + string(filepath.Separator) + info.LibVersion
	if !utils.IsDir(save) {
		if e := utils.MkDirAll(save); e != nil {
			log.Println("创建目录失败 " + save)
			return
		}
	}
	//要保存的文件路径
	saveFile := save + string(filepath.Separator) + info.FileName
	if utils.IsFile(saveFile) {
		//存在的情况下就不要处理了
		log.Println("文件已经存在 " + saveFile)
		return
	}

	req, err := http.NewRequest("GET", info.DownloadURL, nil)
	if err != nil {
		logChan <- fmt.Sprintf("请求失败error :%v\n", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logChan <- fmt.Sprintf("error download failed. %v\n", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		logChan <- fmt.Sprintf("download error : %v\n", resp.StatusCode)
		return
	}

	defer resp.Body.Close()

	//增加一个可视化下载功能
	size := resp.ContentLength
	barLength := 50
	progess := NewBar(info.LibName, barLength, size)
	body := io.TeeReader(resp.Body, progess)

	fp, err := os.OpenFile(saveFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Printf("create file error: %v\n ", err)
		return
	}
	//
	defer fp.Close()

	cnt, e := io.Copy(fp, body)

	if e != nil {
		log.Printf("save file content error : %v\n", e)
		return
	}

	// 同时保存一下描述文件
	libXml := save + string(filepath.Separator) + "library.xml"
	descrition, err := os.OpenFile(libXml, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println("写入xml文件失败  error", err)
		return
	}
	defer descrition.Close()

	var libxmlObj = LibXmlTemplate{
		LibName:        info.LibName,
		LibCategory:    info.Category,
		LibVersion:     info.LibVersion,
		LibComment:     info.Comment,
		LibReleaseDate: info.ReleaseDate,
	}
	//保存一下xml描述在指定的目录下
	xmlContent, err := xml.MarshalIndent(libxmlObj, "", " ")
	if err != nil {
		logChan <- fmt.Sprintf("encode xml内容erro : %v libname: %v\n", err, info.LibName)
		return
	}
	descrition.Write([]byte(xml.Header))
	descrition.Write(xmlContent)
	descrition.Write([]byte("\n"))
	logChan <- fmt.Sprintf("download  %s version: %s writeBytes:%v success\n", info.LibName, info.LibVersion, cnt)

}

func PrintLog() {
	green := "\x1B[32m"
	red := "\x1B[30m"
	reset := "\x1B[0m"
	for logInfo := range logChan {
		color := green
		if strings.Contains(logInfo, "error") {
			color = red
		}
		log.Printf(color+"-> %s\n"+reset, logInfo)
	}
}
