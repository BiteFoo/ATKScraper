package main

import (
	"encoding/xml"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BiteFoo/android_sdk_scraper/download"
	"github.com/BiteFoo/android_sdk_scraper/repos"
	"github.com/BiteFoo/android_sdk_scraper/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoadLibrary(t *testing.T) {
	//
	mvn := "amazon" //"mvn-central"
	result := repos.LoadLibraySpec(mvn)
	log.Printf("=> libraries : %v\n", result.Libraries)
}

func TestRepos(t *testing.T) {

	// mvn := repos.NewMvn()
	// GetReposInfo()

}

// func test

func TestDownload(t *testing.T) {

	utils.MakeSurePathExists("TestingSdk", "")
}

func TestScanZeroByteFile(t *testing.T) {

	pwd, _ := os.Getwd()
	root := pwd + string(filepath.Separator) + "download-lib-repos"

	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.Size() <= 0 {
			log.Println("Found zero bytes file " + path)
			os.Remove(path)
		}

		return err
	})
	//

}

func TestFmtTime(t *testing.T) {

	stamp := 1681107307000 / 1000
	// log.Println(stamp)
	now := time.Now().Unix()

	s1 := time.Unix(int64(stamp), 0).Format("02.01.2006")
	s2 := time.Unix(now, 0).Format("02.01.2006")
	log.Println(stamp, " s1  = ", s1, " ", now, " s2 = ", s2)
}

func TestXml(t *testing.T) {
	lib := download.LibXmlTemplate{

		LibName:        "AndroidAsync",
		LibCategory:    "Android",
		LibVersion:     "1.0.1",
		LibReleaseDate: "06.06.2023",
	}
	// enc := xml.NewEncoder(os.Stdout)
	// enc.Indent(" ", " ")
	// if err := enc.Encode(v); err != nil {
	// 	log.Panicf("encode xml error %v\n", err)
	// }
	output, err := xml.MarshalIndent(lib, "", " ")
	if err != nil {
		log.Panicf("encode xml error %v\n", err)
	}
	// var result []byte
	// result = append(result, []byte(xml.Header))
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
	os.Stdout.Write([]byte("\n"))
}

func TestGoogleMvn(t *testing.T) {

	p := "D:\\CodeDocument\\FridaHookCode\\testapp\\20200912新app\\maven-metadata.xml"
	data, err := ioutil.ReadFile(p)
	assert.Nil(t, err)
	result := repos.MetaDataXML{}

	if e := xml.Unmarshal(data, &result); e != nil {
		log.Panic(e)
	}
	log.Printf("groupId =%v\n", result.GroupId)
	log.Printf("artifactId =%v\n", result.ArteFaceId)
	log.Printf("versioning =%v\n", result.VersioningData)

	log.Printf("versioning.latest =%v\n", result.VersioningData.LastUpdated)
	log.Printf("versioning.release =%v\n", result.VersioningData.RelaseVersion)
	log.Printf("versions =%v\n", result.VersioningData.Versions)
	log.Printf("versioning.lastUpdated =%v\n", result.VersioningData.LastUpdated)
	// log.Printf("groupId =%v\n", result)

}
