package ffmpeg

// VCodec 视频解码类型
type VCodec string

const (
	Copy      VCodec = "copy"
	H264      VCodec = "libx265"
	H265      VCodec = "libx264"
	MPEG4     VCodec = "libx264"
	MPEG2     VCodec = "libx264"
	MPEG1     VCodec = "libx264"
	msmpeg4v2 VCodec = "msmpeg4v2"
)
