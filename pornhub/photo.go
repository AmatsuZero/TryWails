package pornhub

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
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

func (p *Photo) GetPhotos(config DownloadConfig) {
	resp, err := p.loadAlbumPage(config.Page)
	if err != nil {

	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}

	for _, ambumURL := range p.scrapAlbumsURL(doc) {

	}
}

func (p *Photo) loadAlbumPage(pageNum int) (*http.Response, error) {
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

	req, err := getRequest(BASE_URL+ALBUMS_URL+strings.Join(categories, "-"), payload)
	if err != nil {
		return nil, err
	}
	return p.client.Do(req)
}

func (p *Photo) scrapAlbumsURL(doc *goquery.Document) (albumsURL []string) {
	doc.Find("div").
		ChildrenFiltered(".photoAlbumListBlock").
		Each(func(i int, selection *goquery.Selection) {
			selection.Find("a").
				EachWithBreak(func(i int, selection *goquery.Selection) bool {
					u, ok := selection.Attr("href")
					if !ok || !isAlbum(u) {
						return true
					}
					albumsURL = append(albumsURL, BASE_URL+u)
					return false
				})
		})
	return
}

func (p *Photo) scrapPhotoFullURL(previewURL string) (imgURL string, err error) {
	req, err := getRequest(previewURL, nil)
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

func (p *Photo) scrapeAlbumPhotos(albumURL string) {

}
