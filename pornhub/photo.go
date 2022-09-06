package pornhub

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
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

	u, err := url.Parse(BASE_URL + ALBUMS_URL + strings.Join(categories, "-"))
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
	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}

	return p.client.Do(req)
}

func (p *Photo) scrapAlbumsURL(doc *goquery.Document) (albumsURL []string) {
	doc.Find("div").
		ChildrenFiltered(".photoAlbumListBlock").
		Each(func(i int, selection *goquery.Selection) {
			u, ok := selection.Find("a").Attr("href")
			if !ok || !isAlbum(u) {
				return
			}
			albumsURL = append(albumsURL, BASE_URL+u)
		})
	return
}

func (p *Photo) scrapPhotoFullURL(previewURL string) {
	req, err := http.NewRequest("GET", previewURL, nil)
	for k, v := range HEADERS {
		req.Header.Set(k, v)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	src, ok := doc.Find("img").
		FilterFunction(func(i int, selection *goquery.Selection) bool {
			u, ok := selection.Attr("src")
			return ok && isPhoto(u)
		}).First().Attr("src")

}
