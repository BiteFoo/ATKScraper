package repos

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"github.com/BiteFoo/android_sdk_scraper/download"
	"github.com/BiteFoo/android_sdk_scraper/http"
	"github.com/BiteFoo/android_sdk_scraper/utils"
)

//下载google仓库的数据
/**
<?xml version="1.0" encoding="UTF-8"?>
<metadata>
  <groupId>androidx.palette</groupId>
  <artifactId>palette</artifactId>
  <versioning>
    <latest>1.0.0</latest>
    <release>1.0.0</release>
    <versions>
      <version>1.0.0-alpha1</version>
      <version>1.0.0-alpha3</version>
      <version>1.0.0-beta01</version>
      <version>1.0.0-rc01</version>
      <version>1.0.0-rc02</version>
      <version>1.0.0</version>
    </versions>
    <lastUpdated>20180507111125</lastUpdated>
  </versioning>
</metadata>

**/
type MavenMetaDataXML struct {
	Meta MetaDataXML `xml:"metadata"`
}
type MetaDataXML struct {
	GroupId        string     `xml:"groupId"`
	ArteFaceId     string     `xml:"artifactId"`
	VersioningData Versioning `xml:"versioning"`
}

type Versioning struct {
	LatestVersion string  `xml:"latest"`
	RelaseVersion string  `xml:"release"`
	Versions      Version `xml:"versions"`
	LastUpdated   string  `xml:"lastUpdated"`
}

type Version struct {
	Value []string `xml:"version"`
}

const (
	DEFUALT_FILE_TYPE = "aar"
)

// 解析出xml
func pareMetaXML(data []byte) (MetaDataXML, error) {

	var result MetaDataXML
	if er := xml.Unmarshal(data, &result); er != nil {
		return MetaDataXML{}, er
	}

	return result, nil

}

// 下载xml返回结果的值
func DownloadXMLLibs(client *download.DownloadClient, lib RepoLibrary) error {

	// google,trackers.json 内的我们默认下载aar即可 ，不用jar
	//lib内的每个item都是包含了repo:http://xx的值

	for _, repo := range lib.Libraries {
		// 首先需要下载metadata.xml文件
		// log.Println(repo.GroupId, " ", repo.ArteFaceId)
		gr := strings.Replace(repo.GroupId, ".", "/", -1)
		ar := strings.Replace(repo.ArteFaceId, ".", "/", -1)
		metaurl := repo.RepoUrl + "/" + gr + "/" + ar + "/maven-metadata.xml"
		//解析出来xml内容
		data, err := http.Get(metaurl)
		if err != nil {
			log.Println("下载 "+metaurl+" 失败 ", err)
			continue
		}

		// er := io.Copy(dst Writer, src Reader)(data, resp.Body)
		metadata, er := pareMetaXML(data)
		if er != nil {
			log.Println("解析xml失败 => ", er)
			continue
		}
		// log.Println("download xml = ", metadata)
		// 处理所有的version并进行下载
		for _, version := range metadata.VersioningData.Versions.Value {
			fileNmae := fmt.Sprintf("%s-%v.%v", ar, version, DEFUALT_FILE_TYPE) //artifacedIdR + "-" + ".version" + ".filetyp"
			downloadURL := repo.RepoUrl + "/" + gr + "/" + ar + "/" + version + "/" + fileNmae
			// 这里我们就能保证路径创建为 /download/to/category/libname/version
			//
			root := utils.MakeSurePathExists(repo.Category, strings.Replace(repo.Name, "::", "_", -1))
			if root == "" {
				return fmt.Errorf("无法创建目录")
			}
			//我们在这里将结果保存在chan里提供给别的goroutine下载
			var info = download.DownloadInfo{
				Category:    repo.Category,
				DownloadURL: downloadURL,
				LibName:     repo.Name,
				LibVersion:  version,
				FileName:    fileNmae,
				SaveRoot:    root,
				Comment:     repo.Commment,
				ReleaseDate: metadata.VersioningData.LastUpdated,
			}
			// 提供给download执行
			client.Submit(info)
			log.Printf("-> metaXml downloadInfo: %v\n", info)
			// return nil

		}

	}
	return nil

}
