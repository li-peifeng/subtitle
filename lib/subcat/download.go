package subcat

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func OpenDownloadPage(path string, downloadLangCode string) string {
	u := buildURL(baseURL, path, nil)
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	selector := "#download_" + downloadLangCode
	selection := doc.Find(selector)
	if selection.Length() == 0 {
		log.Fatal("No subtitle found for language:", downloadLangCode)
	}

	href, exists := selection.Attr("href")
	if !exists {
		log.Fatalf("%s has no href attribute found for language %s", u, downloadLangCode)
	}

	href, _ = url.QueryUnescape(href)

	return href
}

func DownloadSubFile(path string) []byte {
	url := buildURL(baseURL, path, nil)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	encoding, _, _ := charset.DetermineEncoding(body, "")

	if encoding == unicode.UTF8 || encoding == nil {
		return body
	}

	utf8Reader := transform.NewReader(bytes.NewReader(body), encoding.NewDecoder())
	utf8Body, err := io.ReadAll(utf8Reader)
	if err != nil {
		log.Fatalf("Error converting to UTF-8: %v", err)
	}

	return utf8Body
}
