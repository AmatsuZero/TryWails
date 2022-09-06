package pornhub

import (
	"net/http"
	"path/filepath"
	"strings"
)

const (
	BASE_URL   = "https://pornhub.com"
	PHOTO_EXT  = ".jpg" // for validation
	SEARCH_URL = "/search"

	PORNSTARS_URL  = "/pornstars"
	PORNSTAR_URL   = "/pornstar/" // for validation
	MODEL_URL      = "/model/"
	PORNSTAR_PHOTO = ".phncdn.com/" // for validation

	VIDEOS_URL      = "/video"
	VIDEO_URL       = "/view_video.php?viewkey=" // for validation
	VIDEO_IMAGE_URL = ".phncdn.com/videos/"      // for validation

	ALBUMS_URL      = "/albums/"
	ALBUM_URL       = "/album/"                 // for validation
	ALBUM_PHOTO_URL = "phncdn.com/pics/albums/" // for validation
	PHOTO_PREVIEW   = "/photo/"                 // for validation

	TIME_TO_WAIT = 3
)

var HEADERS map[string]string

func init() {
	HEADERS = map[string]string{
		"Content-Type": "text/html; charset=UTF-8",
	}
}

func isAlbum(url string) bool {
	return strings.Contains(url, ALBUMS_URL)
}

func isPhotoPreview(url string) bool {
	return strings.Contains(url, PHOTO_PREVIEW)
}

func isPhoto(url string) bool {
	return strings.Contains(url, ALBUM_PHOTO_URL) && filepath.Ext(url) == PHOTO_EXT
}

func isStar(url string) bool {
	return strings.Contains(url, PORNSTAR_URL) || strings.Contains(url, MODEL_URL)
}

func isStarPhoto(url string) bool {
	return strings.Contains(url, PORNSTAR_PHOTO) && filepath.Ext(url) == PHOTO_EXT
}

func isVideo(url string) bool {
	return strings.Contains(url, VIDEO_URL)
}

func isVideoPhoto(url string) bool {
	return strings.Contains(url, VIDEO_IMAGE_URL) && filepath.Ext(url) == PHOTO_EXT
}

type DownloadConfig struct {
	Quantity, Page int
	InfinityRetry  bool
}

type PornHub struct {
	client *http.Client
	Photos *Photo
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
