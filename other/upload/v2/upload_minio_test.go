package upload

import (
	"github.com/injoyai/goutil/cache/v2"
	"github.com/injoyai/goutil/oss"
	"os"
	"testing"
)

func TestNewMinio(t *testing.T) {
	cfg := cache.NewFile(oss.UserInjoyDir("/data/cache/cmd"), "global")
	up, err := NewMinio(&MinioConfig{
		Endpoint:   cfg.GetString("uploadMinio.endpoint"),
		AccessKey:  cfg.GetString("uploadMinio.access"),
		SecretKey:  cfg.GetString("uploadMinio.secret"),
		BucketName: cfg.GetString("uploadMinio.bucket"),
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	f, err := os.Open("./upload_minio.go")
	if err != nil {
		t.Fatal(err)
		return
	}
	defer f.Close()

	u, err := up.Upload("dir/upload_minio.go", f)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(u)

	ls, err := up.List("dir/")
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, v := range ls {
		t.Log(*v)
	}
}
