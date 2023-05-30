package in

import (
	"github.com/injoyai/base/oss"
	"github.com/injoyai/conv"
	"net/http"
)

//Read 解析json数据
func Read(r interface{}, ptr interface{}) {
	read(r, ptr)
}

//ReadJson 解析json数据
func ReadJson(r interface{}, ptr interface{}) {
	Read(r, ptr)
}

func GetBodyString(r interface{}) string {
	return string(getBody(r))
}

func GetBody(r interface{}) []byte {
	return getBody(r)
}

func Header(r interface{}) http.Header {
	return getHeader(r)
}

func GetHeader(r interface{}, key string) string {
	return getHeader(r).Get(key)
}

func GetFile(r interface{}, name string) []byte {
	return getFile(r, name)
}

func SaveFile(r interface{}, name, path string) {
	CheckErr(oss.New(path, getFile(r, name)))
}

func GetPageSize(r interface{}) (int, int) {
	return getPageSize(r)
}

func Get(r interface{}, key string, def ...interface{}) *conv.Var {
	return get(r, key, def...)
}

func GetString(r interface{}, key string, def ...interface{}) string {
	return get(r, key, def...).String()
}

func GetInt(r interface{}, key string, def ...interface{}) int {
	return get(r, key, def...).Int()
}

func GetInt64(r interface{}, key string, def ...interface{}) int64 {
	return get(r, key, def...).Int64()
}

//func GetLocalFile(path string) (interface{}, error) {
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
