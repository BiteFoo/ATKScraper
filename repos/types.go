package repos

import "github.com/BiteFoo/android_sdk_scraper/download"

// 仓库包的数据
// 从配置文件中读取要爬取的库信息
type RepoLibrary struct {
	Libraries []RepoInfo `json:"libraries"`
}

/*
*
libray-specs
repo struct

	"artefactid": "preference-ktx",
	"category": "Android",
	"comment": "",
	"groupid": "androidx.preference",
	"name": "androidx.preference::preference-ktx",
	"repo": "https://dl.google.com/dl/android/maven2"
*
*/
type RepoInfo struct {
	ArteFaceId string `json:"artefactid"`
	Category   string `json:"category"`
	Commment   string `json:"comment"`
	GroupId    string `json:"groupid"`
	Name       string `json:"name"`
	RepoUrl    string `json:"repo"`
}

type Repo struct {
	repoMap map[string]RepoLibrary
	Client  *download.DownloadClient
	// BaseUrl string
}

// 下载保存信息
type SaveInfo struct {
	SavePath    string //保存路径 //
	ArteFaceId  string // 对应sdk名
	GroupId     string //
	FileContent []byte //文件内容
	Type        string //标识分类用
	XmlContent  []byte // 保存的xml文件内容 默认aar/jar保存在同一个目录下，命名为 library.xml
}
