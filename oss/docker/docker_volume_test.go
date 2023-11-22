package docker

import (
	"testing"
)

func TestClient_VolumeList(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Error(err)
		return
	}
	list, _, err := c.VolumeList(&VolumeSearch{
		PageSize: PageSize{},
		Name:     "",
		Driver:   "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range list {
		t.Log(*v)
	}
}
