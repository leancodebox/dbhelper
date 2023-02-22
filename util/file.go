package util

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/purerun/dbhelper/util/eh"

	"github.com/spf13/cast"
)

var basePath string

func SetBasePath(path string) {
	basePath = path
}

func GetStoragePath(filename string) string {
	return basePath + filename
}

func StorageGet(filename string) string {
	data, _ := FileGetContents(basePath + filename)
	return cast.ToString(data)
}

func StoragePut(filename string, data any, append bool) error {
	return FilePutContents(basePath+filename, data, append)
}

// Put 将数据存入文件
func Put(data []byte, to string) (err error) {
	err = os.WriteFile(to, data, 0644)
	return
}

func FileGetContents(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// FilePutContents file_put_contents
func FilePutContents(filename string, data any, isAppend ...bool) error {
	if dir := filepath.Dir(filename); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	bData := []byte(cast.ToString(data))
	needAppend := false
	if len(isAppend) > 0 && isAppend[0] == true {
		needAppend = true
	}
	if needAppend {
		fl, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer func(fl *os.File) {
			eh.PrIF(fl.Close())
		}(fl)
		_, err = fl.Write(bData)
		return err
	} else {
		return os.WriteFile(filename, bData, 0644)
	}
}

func PutContent(filename string, data any) {
	_ = FilePutContents(filename, data, false)
}

func AppendPutContent(filename string, data any) {
	_ = FilePutContents(filename, data, true)
}

func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsExistOrCreate(path string, init string) bool {
	if IsExist(path) {
		return true
	}
	PutContent(path, init)
	return true
}

func DirExistOrCreate(dirPath string) bool {
	if IsExist(dirPath) {
		return true
	} else {
		return os.MkdirAll(dirPath, os.ModePerm) != nil
	}
}

func UrlDecode(s string) string {
	r, err := url.QueryUnescape(s)
	if err != nil {
		return ""
	}
	return r
}

func UrlEncode(s string) string {
	return url.QueryEscape(s)
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
