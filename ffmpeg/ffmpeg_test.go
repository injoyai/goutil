package ffmpeg

import "testing"

func TestCapture(t *testing.T) {
	t.Log(Capture("F:\\test\\prog_index.ts", "00:00:37", "F:\\test\\test.jpg"))
}

func TestGif(t *testing.T) {
	DefaultDebug = true
	t.Log(Gif("F:\\test\\prog_index.ts", "00:00:37", 4, 30, 320, 240, "F:\\test\\test.gif"))
}

func TestToAudio(t *testing.T) {
	DefaultDebug = true
	t.Log(ToAudio("F:\\test\\prog_index.ts", "37", 10, "F:\\test\\test.mp3"))
}
