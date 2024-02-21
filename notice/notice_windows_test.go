package notice

import (
	"testing"
)

func TestNewWindows(t *testing.T) {
	t.Log(NewWindows().Publish(&Message{
		Title:   "标题",
		Content: "内容",
	}))
	t.Log(NewWindows().Publish(&Message{
		Target:  TargetPopup,
		Title:   "标题",
		Content: "内容",
	}))
}
