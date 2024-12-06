### 说明
* 对[`ffmpeg`](`https://ffmpeg.org/`)的使用封装
* 下载地址[`latest`](`https://github.com/BtbN/FFmpeg-Builds/releases`)

### 教程


* 从一个视频中，在10s位置开始截取，截取12s，输出格式为flv格式，输出名字为out.mp4
    ```
    ffmpeg -i xxx.mp4 -ss 10 -t 12 -f flv out.mp4
    ```

* 在一个视频中提取音频，输出400帧音频，音频编码率为192k，采样率为48k，通道数为2，音频编码为libmp3lame，输出名字为out.mp3
    ```
    ffmpeg -i out.mp4 -aframes 400 -b:a 192k -ar 48000 -ac 2 -acodec libmp3lame out.mp3
    ```
* 从一个mp4视频中，输出视频400帧数，视频编码为300k，帧速率为60，画面大小为640x480，纵横比为16:9，视频编码为libx265，输出名字为out.h265
  ```
  ffmpeg -i xxx.mp4 -vframes 400 -b:v 300k -r 60 -s 640x480 -aspect 16:9 -acodec copy -vcodec libx265 out.h265
  ```

#### 参考
* 链接[`百度`](`https://baike.baidu.com/item/ffmpeg/2665727?fr=ge_ala`)
* 链接[`CSDN`](`https://blog.csdn.net/cpp_learner/article/details/142657414`)
