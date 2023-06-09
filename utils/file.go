package utils

import (
	"log"
	"os"
	"path/filepath"
)

// 获取当前目录
func GetPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		// 尝试获取运行的二进制路径
		pwd, err = os.Executable()
	}

	if err != nil {
		log.Println("读取pwd失败 ", err)
		return ""
	}

	return pwd

}

func IsFile(p string) bool {
	st, err := os.Stat(p)
	if err != nil {
		//
		// log.Printf("读取文件stat 失败 error: %v\n", err)
		return false
	}
	return !st.IsDir()
}

func IsDir(p string) bool {
	st, err := os.Stat(p)
	if err != nil {
		return false
	}
	return st.IsDir()
}

func MkDirAll(p string) error {
	return os.Mkdir(p, 0755)

}

// 创建目录保存 失败返回空
func MakeSurePathExists(category, libName string) string {

	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return ""
	}
	// 创建子目录保存
	save := root + string(filepath.Separator) + "download-lib-repos" + string(filepath.Separator) + category + string(filepath.Separator) + libName
	if !IsDir(save) {
		if err := os.MkdirAll(save, 0755); err != nil {
			log.Println(err)
			return ""
		}
	}

	return save

}

func GenLibSpecPath(name string) string {
	// -- *** --
	// -- *** --
	// -- *** --
	return GetPwd() + string(filepath.Separator) + "library-specs" + string(filepath.Separator) + name
}
