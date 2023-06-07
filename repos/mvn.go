package repos

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BiteFoo/android_sdk_scraper/download"
	"github.com/BiteFoo/android_sdk_scraper/http"
	"github.com/BiteFoo/android_sdk_scraper/utils"
)

type MvnRepo struct {
	ReponseHeader MvnResponseHeader `json:"responseHeader"`
	Response      MvnResponse       `json:"response"`
}

type MvnResponse struct {
	NunFound int      `json:"numFound"`
	Start    int      `json:"start"`
	Docs     []MvnDoc `json:"docs"`
}

type MvnDoc struct {
	Id         string   `json:"id"`
	Group      string   `json:"g"`
	ArteFaceId string   `json:"a"`
	Version    string   `json:"v"`
	P          string   `json:"p"`
	Timestamp  int64    `json:"timestamp"`
	Ec         []string `json:"ec"`
	Tags       []string `json:"tags"`
}
type MvnResponseHeader struct {
	Status int        `json:"status"`
	QTime  int        `json:"QTime"`
	Params RespParams `json:"params"`
}

type RespParams struct {
	Q       string `json:"q"`
	Core    string `json:"core"`
	Indent  string `json:"indent"`
	Fl      string `json:"fl"`
	Start   string `json:"start"`
	Sort    string `json:"sort"`
	Row     string `json:"row"`
	Wt      string `json:"wt"`
	Version string `json:"version"`
}

// maven updater

func NewMvn() *Repo {

	libs := LoadLibraySpec("mvn-central")

	client := download.NewDownloader()
	return &Repo{
		Repos:   libs,
		BaseUrl: "http://search.maven.org/solrsearch/select?q=g:%22#groupId%22+AND+a:%22#artefactId%22&rows=100&core=gav",
		Client:  client,
	}
}

// 读取指定的libray-specs数据
func LoadLibraySpec(name string) *RepoLibrary {
	//加载出指定的lib
	libPath := utils.GenLibSpecPath(name)
	if !utils.IsFile(libPath) {
		return nil
	}
	var repoLibs RepoLibrary
	r, _ := os.OpenFile(libPath, os.O_RDONLY, os.ModePerm)

	err := json.NewDecoder(r).Decode(&repoLibs)
	if err != nil {
		log.Printf("Read library-spec file err : %v\n", err)
		return nil
	}
	return &repoLibs

}

func formatString(url string, repo RepoInfo) string {

	url = strings.Replace(url, "#groupId", repo.GroupId, -1)
	url = strings.Replace(url, "#artefactId", repo.ArteFaceId, -1)
	return url
}

func GetReposInfo(r *Repo) error {

	for _, repo := range r.Repos.Libraries {

		url := formatString(r.BaseUrl, repo)
		// log.Printf("Fetching repos[%d] url = %v\n", i, url)
		result, err := http.Get(url)
		if err != nil {
			log.Printf("http errro : %v\n", err)
			continue
		}
		// log.Printf("size(result) = %v %v\n", len(result), string(result))
		var resp MvnRepo
		if er := json.Unmarshal(result, &resp); er != nil {
			log.Printf("decode response error; %v\n", er)
			continue
		}
		downloadArtefaceId(r, repo, resp)
		// break

	}
	// close(r.Client.DownloadChans) //不要开这个，否则会有协程没有下载完成？
	return nil

}

func parseTime(sec int64) string {
	stamp := sec / 1000 // 去掉尾部的000
	s1 := time.Unix(int64(stamp), 0).Format("02.01.2006")
	return s1
}

// 下载 获取响应结果保存结果
func downloadArtefaceId(r *Repo, repo RepoInfo, resp MvnRepo) error {

	artifacedIdR := strings.Replace(repo.ArteFaceId, ".", "/", -1)
	groupIdR := strings.Replace(repo.GroupId, ".", "/", -1)

	baseUrl := "https://search.maven.org/remotecontent?filepath="

	for _, respons := range resp.Response.Docs {
		//  response.P == response.fileType
		fileNmae := fmt.Sprintf("%s-%v.%v", artifacedIdR, respons.Version, respons.P) //artifacedIdR + "-" + ".version" + ".filetyp"
		downloadURL := baseUrl + groupIdR + "/" + artifacedIdR + "/" + respons.Version + "/" + fileNmae
		// log.Println("[", i, "]fileName = ", fileNmae, " DownloadURL = ", downloadURL)
		// break
		// 这里我们就能保证路径创建为 /download/to/category/libname/version
		root := utils.MakeSurePathExists(repo.Category, repo.Name)
		if root == "" {
			return fmt.Errorf("无法创建目录")
		}
		//

		//我们在这里将结果保存在chan里提供给别的goroutine下载

		var info = download.DownloadInfo{
			Category:    repo.Category,
			DownloadURL: downloadURL,
			LibName:     repo.Name,
			LibVersion:  respons.Version,
			FileName:    fileNmae,
			SaveRoot:    root,
			Comment:     repo.Commment,
			ReleaseDate: parseTime(respons.Timestamp),
		}
		// 提供给download执行
		r.Client.Submit(info)
	}

	return nil
}
