package http_handler

import (
	"net/http"
	"testing"
)

func TestNewFile(t *testing.T) {
	t.Log(http.ListenAndServe(":8200", DefaultFile))
}

func TestNewFile2(t *testing.T) {
	t.Log(http.ListenAndServe(":8200", NewFile("", "/api")))
}
