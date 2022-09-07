package pornhub

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

type Photo struct {
	keywords []string
	client   *http.Client
}

func newPhoto(c *http.Client, kw []string) *Photo {
	return &Photo{
		keywords: kw,
		client:   c,
	}
}

func (p *Photo) GetPhotos(config DownloadConfig, proc func(url string, err error) bool) {
	ctx, cancel := context.WithCancel(context.Background())
	quantity := uint64(config.Quantity)

	if quantity < 1 {
		quantity = 1
	}

	var found uint64 = 0
	queue := NewQueue()

	quit := func() {
		cancel()
		queue.Stop()
	}

	page := config.Page
	if page < 1 {
		page = 1
	}

	for executed := false; !executed || config.Infinity; page++ {
		executed = true

		resp, err := p.loadAlbumPage(page, ctx)
		if err != nil {
			if proc("", err) {
				quit()
				return
			}
			continue
		}

		urls, err := p.scrapAlbumsURL(resp)
		if err != nil {
			if proc("", err) {
				quit()
				return
			}
			continue
		}

		for _, albumURL := range urls {
			queue.Invoke(func() {
				u, err := p.scrapAlbumPhotos(albumURL, ctx)
				if err != nil {
					if proc(u, err) {
						quit()
					}
					return
				}

				u, err = p.scrapPhotoFullURL(u, ctx)
				if err != nil {
					if proc(u, err) {
						quit()
					}
					return
				}

				shouldCancel := proc(u, err)
				atomic.AddUint64(&found, 1)
				if shouldCancel || found == quantity {
					quit()
					return
				}
			})
		}
	}
}

func (p *Photo) loadAlbumPage(pageNum int, ctx context.Context) (*http.Response, error) {
	payload := map[string]string{
		"search": "",
		"page":   strconv.Itoa(pageNum),
	}
	var categories []string
	var searchWords []string
	for _, keyword := range p.keywords {
		switch keyword {
		case "female", "straight", "misc", "male", "gay":
			categories = append(categories, keyword)
		default:
			searchWords = append(searchWords, keyword)
		}
	}
	if len(searchWords) > 0 {
		payload["search"] = strings.Join(searchWords, "+")
	}

	req, err := getRequest(BaseUrl+AlbumsUrl+strings.Join(categories, "-"), payload, ctx)
	if err != nil {
		return nil, err
	}
	return p.client.Do(req)
}

func (p *Photo) scrapAlbumsURL(resp *http.Response) (albumsURL []string, err error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	doc.Find("div").Filter(".photoAlbumListBlock").
		Each(func(i int, selection *goquery.Selection) {
			selection.Find("a").
				EachWithBreak(func(i int, selection *goquery.Selection) bool {
					u, ok := selection.Attr("href")
					if !ok || !isAlbum(u) {
						return true
					}
					albumsURL = append(albumsURL, BaseUrl+u)
					return false
				})
		})
	return
}

func (p *Photo) scrapPhotoFullURL(previewURL string, ctx context.Context) (imgURL string, err error) {
	req, err := getRequest(previewURL, nil, ctx)
	if err != nil {
		return "", err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	doc.Find("img").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		u, ok := selection.Attr("src")
		if !ok || !isPhoto(u) {
			return true
		}
		imgURL = u
		return false
	})
	return imgURL, err
}

func (p *Photo) scrapAlbumPhotos(albumURL string, ctx context.Context) (string, error) {
	r, err := getRequest(albumURL, nil, ctx)
	if err != nil {
		return "", err
	}

	resp, err := p.client.Do(r)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	u := ""
	doc.Find("a").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		a, ok := selection.Attr("href")
		if !ok || !isPhotoPreview(a) {
			return true
		}
		u = BaseUrl + a
		return false
	})
	return u, nil
}
