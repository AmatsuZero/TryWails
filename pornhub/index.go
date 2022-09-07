package pornhub

import (
	"context"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	BaseUrl    = "https://pornhub.com"
	PhotoExt   = ".jpg" // for validation
	SEARCH_URL = "/search"

	PORNSTARS_URL = "/pornstars"
	PornStarUrl   = "/pornstar/" // for validation
	ModelUrl      = "/model/"
	PornStarPhoto = ".phncdn.com/" // for validation

	VIDEOS_URL    = "/video"
	VideoUrl      = "/view_video.php?viewkey=" // for validation
	VideoImageUrl = ".phncdn.com/videos/"      // for validation

	AlbumsUrl     = "/albums/"
	AlbumUrl      = "/album/"                 // for validation
	AlbumPhotoUrl = "phncdn.com/pics/albums/" // for validation
	PhotoPreview  = "/photo/"                 // for validation

	TIME_TO_WAIT = 3
)

var HEADERS map[string]string

func init() {
	HEADERS = map[string]string{
		"Content-Type": "text/html; charset=UTF-8",
	}
}

func isAlbum(url string) bool {
	return strings.Contains(url, AlbumUrl)
}

func isPhotoPreview(url string) bool {
	return strings.Contains(url, PhotoPreview)
}

func isPhoto(url string) bool {
	return strings.Contains(url, AlbumPhotoUrl) && filepath.Ext(url) == PhotoExt
}

func isStar(url string) bool {
	return strings.Contains(url, PornStarUrl) || strings.Contains(url, ModelUrl)
}

func isStarPhoto(url string) bool {
	return strings.Contains(url, PornStarPhoto) && filepath.Ext(url) == PhotoExt
}

func isVideo(url string) bool {
	return strings.Contains(url, VideoUrl)
}

func isVideoPhoto(url string) bool {
	return strings.Contains(url, VideoImageUrl) && filepath.Ext(url) == PhotoExt
}

func getRequest(str string, payload map[string]string, ctx context.Context) (*http.Request, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for k, v := range payload {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}

	return req, err
}

type DownloadConfig struct {
	Quantity, Page int
	Infinity       bool
}

type PornHub struct {
	client *http.Client
	Photos *Photo
}

func NewQueue() *DispatchQueue {
	q := &DispatchQueue{work: make(chan func()), quit: make(chan bool), Working: true}
	go func() {
		var job func()
		for {
			select {
			case job = <-q.work:
			case <-q.quit:
				q.Working = false
				return

			}
			job()
		}
	}()
	return q
}

func (q *DispatchQueue) Invoke(work func()) {
	q.work <- work
}

func (q *DispatchQueue) Stop() {
	if !q.Working {
		return
	}
	q.quit <- true
}

func (q *DispatchQueue) Start() {
	if q.Working {
		return
	}

	go func() {
		var job func()

		q.Working = true

		for {
			select {
			case job = <-q.work:
			case <-q.quit:
				q.Working = false
				return

			}
			job()
		}
	}()
}

type DispatchQueue struct {
	work    chan func() // 任务队列
	quit    chan bool   // 退出标志
	Working bool        // 运行标志
}

func NewPornHub(keywords []string) *PornHub {
	client := &http.Client{
		Transport: &http.Transport{},
	}
	return &PornHub{
		client: client,
		Photos: newPhoto(client, keywords),
	}
}
