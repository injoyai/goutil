package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/oss"
	"github.com/injoyai/logs"
	"github.com/injoyai/lorca"
)

type APP = lorca.APP

//go:embed config.html
var html string

// GUI 界面 未完成
func GUI(cfg *Config) error {
	if cfg.Width <= 0 {
		cfg.Width = 720
	}
	if cfg.Height <= 0 {
		cfg.Height = 480
	}

	// 启动本地 HTTP 服务器
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	})
	mux.HandleFunc("/api/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg.Frames)
	})
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			json.NewDecoder(r.Body).Decode(cfg.Frames)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfg.Frames)
	})

	server := &http.Server{Addr: ":0", Handler: mux}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		logs.PanicErr(err)
	}
	go server.Serve(listener)
	defer server.Close()

	return lorca.Run(&lorca.Config{
		Width:  cfg.Width,
		Height: cfg.Height,
		Index:  fmt.Sprintf("http://127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port),
	}, func(app lorca.APP) error {
		err := cfg.Loading()
		if err != nil {
			logs.Err(err)
		}
		//绑定函数
		app.Bind("SaveConfig", func(jsonStr string) {
			logs.Debug(jsonStr)
			cfg.Save(jsonStr)
		})
		app.Bind("GetSchema", func() string {
			bs, _ := json.Marshal(cfg.Frames)
			return string(bs)
		})
		app.Bind("GetConfig", func() string {
			return cfg.m.String()
		})
		return nil
	})
}

type Config struct {
	Width    int       //宽度,可选
	Height   int       //高度,可选
	Filename string    //本地文件路径,必须
	Frames   []Frame   //格式,必须
	m        *conv.Map //缓存数据
}

func (this *Config) Loading() error {
	bs, err := os.ReadFile(this.Filename)
	if err != nil {
		return err
	}
	this.m = conv.NewMap(bs)
	return nil
}

func (this *Config) Save(jsonStr string) error {
	err := oss.New(this.Filename, jsonStr)
	if err != nil {
		return err
	}
	this.m = conv.NewMap(jsonStr)
	return nil
}
