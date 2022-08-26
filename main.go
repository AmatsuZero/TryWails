package main

import (
	"doodle/nhentai"
	_ "embed"
	"github.com/wailsapp/wails"
	"net/http"
	"os"
	"time"
)

func basic() string {
	return nhentaiClient.PageUrl(2075216, 2, "j")
}

//go:embed frontend/dist/app.js
var js string

//go:embed frontend/dist/app.css
var css string

// nhentail 网站
var nhentaiClient nhentai.Client

func init() {
	nhentaiClient = nhentai.Client{}
	nhentaiClient.Transport = &http.Transport{
		TLSHandshakeTimeout:   time.Second * 10,
		ExpectContinueTimeout: time.Second * 10,
		ResponseHeaderTimeout: time.Second * 10,
		IdleConnTimeout:       time.Second * 10,
	}
}

func main() {
	// 获取 mode
	wails.BuildMode = os.Getenv("BuildMode")
	app := wails.CreateApp(&wails.AppConfig{
		Width:            1024,
		Height:           768,
		Title:            "doodle",
		JS:               js,
		CSS:              css,
		Colour:           "#131313",
		DisableInspector: false,
	})
	app.Bind(basic)
	app.Run()
}
