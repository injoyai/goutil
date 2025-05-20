package in

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"net/http"
)

// Read 解析json数据
func Read(r any, ptr any) {
	read(r, ptr)
}

// ReadJson 解析json数据
func ReadJson(r any, ptr any) {
	Read(r, ptr)
}

func GetBodyString(r any) string {
	return string(getBody(r))
}

func GetBody(r any) []byte {
	return getBody(r)
}

func Header(r any) http.Header {
	return getHeader(r)
}

func GetHeader(r any, key string) string {
	return getHeader(r).Get(key)
}

func GetFile(r any, name string) []byte {
	return getFile(r, name)
}

func SaveFile(r any, name, path string) {
	CheckErr(oss.New(path, getFile(r, name)))
}

func GetPageSize(r any) (int, int) {
	return getPageSize(r)
}

func Get(r any, key string, def ...any) *conv.Var {
	return get(r, key, def...)
}

func GetString(r any, key string, def ...any) string {
	return get(r, key, def...).String()
}

func GetInt(r any, key string, def ...any) int {
	return get(r, key, def...).Int()
}

func GetInt64(r any, key string, def ...any) int64 {
	return get(r, key, def...).Int64()
}

//func GetLocalFile(path string) (any, error) {
//	file, err := os.Open(path)
//	if err != nil {
//		return nil, err
//	}
//	fileInfo, err := file.Stat()
//	if err != nil {
//		return nil, err
//	}
//	if !fileInfo.IsDir() {
//		return file
//	}
//	listFile, err := ioutil.ReadDir(path)
//	in.CheckErr(err)
//	listName := []string{}
//	for _, v := range listFile {
//		listName = append(listName, v.Name())
//	}
//	in.Succ(listName)
//}
