package docker

import (
	"testing"
)

func TestClient_ImageList(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Error(err)
		return
	}
	list, _, err := c.ImageList(&ImageSearch{
		PageSize: PageSize{},
		ID:       "",
		Tag:      "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range list {
		t.Logf("%#v", v)
	}
}
