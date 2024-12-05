package upload

import (
	"github.com/injoyai/goutil/cache/v2"
	"github.com/injoyai/goutil/oss"
	"testing"
)

func TestNewMinio(t *testing.T) {
	cfg := cache.NewFile(oss.UserInjoyDir("/data/cache/cmd"), "global")
	t.Log(cfg.GetString("uploadMinio.endpoint"))
	up, err := NewMinio(&MinioConfig{
		Endpoint:   cfg.GetString("uploadMinio.endpoint"),
		AccessKey:  cfg.GetString("uploadMinio.access"),
		SecretKey:  cfg.GetString("uploadMinio.secret"),
		BucketName: cfg.GetString("uploadMinio.bucket"),
		Rename:     cfg.GetBool("uploadMinio.rename"),
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	ls, err := up.List()
	if err != nil {
		t.Fatal(err)
		return
	}
	for _, v := range ls {
		t.Log(*v)
	}
}
