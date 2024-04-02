package upload

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/oss"
	"testing"
)

func TestNewMinio(t *testing.T) {
	m := conv.NewMap(cache.NewFile(oss.UserInjoyDir("data/cache/cmd"), "global").GMap())
	i, err := NewMinio(&MinioConfig{
		Endpoint:   m.GetString("uploadMinio.endpoint"),
		AccessKey:  m.GetString("uploadMinio.access"),
		SecretKey:  m.GetString("uploadMinio.secret"),
		BucketName: m.GetString("uploadMinio.bucket"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(i.List())

}
