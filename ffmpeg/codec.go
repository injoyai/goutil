package ffmpeg

const (
	Copy   = "copy"   //保存原有编解码
	Concat = "concat" //拼接视频

	H261       = "h261"
	H263       = "h263"
	H263P      = "h263p"
	H264       = "libx264"    //目前最常用的视频编码格式，兼容性好，压缩效率高。
	H265       = "libx265"    //比 H.264 更高效的编码格式，能够在更低的比特率下保持较好的视频质量。
	HEVC       = "libx265"    //比 H.264 更高效的编码格式，能够在更低的比特率下保持较好的视频质量。
	VP8        = "libvpx"     //Google 开发的视频编解码格式，主要用于 WebM。
	VP9        = "libvpx-vp9" //VP8 的继任者，广泛应用于高效视频压缩，支持 4K 视频。
	AV1        = "libaom-av1" //现代的视频编码格式，比 HEVC 更高效，适用于未来的高效视频压缩
	MPEG2      = "mpeg2video" //较旧的标准，广泛用于 DVD 和数字电视。
	MPEG4      = "mpeg4"      //较为常见的视频编码格式，广泛应用于移动设备
	Theora     = "libtheora"  //一种开源的视频编码格式，广泛应用于 WebM 和 Ogg 容器。
	ProRes     = "prores_ks"  //prores_aw //Apple 专用的视频编码格式，主要用于专业视频制作
	DNxHD      = "dnxhd"      //Avid 开发的视频编解码器，主要用于广播和专业影视制作。
	DNxHR      = "dnxhr"      //Avid 开发的视频编解码器，主要用于广播和专业影视制作。
	FlashVideo = "flv"        //Flash Video  用于 Flash 播放器的视频格式。

	AAC        = "aac"        //高效的音频编码格式，广泛应用于流媒体和移动设备。
	MP3        = "libmp3lame" //最流行的音频压缩格式，兼容性好。
	Opus       = "libopus"    //高效的音频编解码器，特别适用于语音通信。
	Vorbis     = "libvorbis"  //开源音频编码格式，通常与 Ogg 容器一起使用。
	FLAC       = "flac"       // 无损音频编码格式，适用于高质量音频存储。
	ALAC       = "alac"       //Apple 的无损音频编码格式。
	AC3        = "ac3"        //Dolby Digital 编码格式，广泛用于 DVD 和蓝光光盘。
	EAC3       = "eac3"       //Dolby Digital Plus 编码格式，适用于现代流媒体和广播。
	WMA        = "wmav2"      //Windows 媒体音频格式，微软的专有音频编解码器。
	Speex      = "speex"      //用于语音压缩的音频编解码器。
	G711       = "g711"       //用于电话语音通信的音频编码格式。
	MPEG2Audio = "mp2"        //适用于 MPEG-2 视频的音频编码格式。

	MP4  = "mp4"      //广泛支持的视频和音频容器，常见于流媒体、移动设备等。
	MKV  = "matroska" //开源的多媒体容器，支持多种视频和音频编码格式。
	AVI  = "avi"      //经典的视频容器，广泛支持不同的视频编解码器。
	WebM = "webm"     //Google 推出的开源容器，通常用于 Web 视频，使用 VP8 或 VP9 视频编码。
	FLV  = "flv"      //Flash 视频容器格式，曾经广泛用于流媒体。
	MOV  = "mov"      //Apple 的视频容器格式，广泛用于专业视频编辑和存储。
	OGG  = "ogg"      //用于容纳 Vorbis 音频和 Theora 视频等开源编解码格式的容器。
	TS   = "mpegts"   //MPEG 传输流，广泛用于广播和流媒体

	SRT = "subrip" //最常见的字幕格式，通常与视频一起使用。
	ASS = "ass"    //高级字幕格式，支持更丰富的排版和样式。
	SSA = "ssa"    //高级字幕格式，支持更丰富的排版和样式。
	VTT = "webvtt" //Web 视频字幕格式，通常用于 HTML5 视频播放器。

	JPEG = "mjpeg" //一种常见的静态图像编码格式，适用于视频流中的帧。
	PNG  = "png"   //一种无损压缩的图像格式，适用于高质量图像
	GIF  = "gif"   //适用于动画的图像格式。
	TIFF = "tiff"  //高质量的无损图像格式，广泛应用于扫描和高分辨率图像

)
