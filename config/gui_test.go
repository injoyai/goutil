package config

import "testing"

func TestGUI(t *testing.T) {
	err := GUI(&Config{
		Width:    800,
		Height:   600,
		Filename: "config.json",
		Frames: []Frame{
			{
				Title: "📋 基础信息",
				Fields: []Field{
					{
						Key:   "name",
						Label: "名称",
						Type:  "input",
					},
				},
			},
		},
	})
	t.Error(err)
}
