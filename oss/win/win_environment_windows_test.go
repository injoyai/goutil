package win

import "testing"

func TestGetRootEnv(t *testing.T) {
	t.Log(GetRootEnv("GOBIN"))
}
