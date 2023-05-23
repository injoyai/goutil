package ip

import (
	"testing"
)

func TestGetLocal(t *testing.T) {
	t.Log(GetLocal())
}

func TestGetLocalAll(t *testing.T) {
	t.Log(GetLocalAll())
}
