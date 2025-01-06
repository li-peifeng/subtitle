package subcat

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

type SearchResult struct {
	Title     string
	Path      string
	Downloads int
}

func Search(keyword string) []SearchResult {
	u := buildURL(baseURL, "index.php", map[string]string{"search": keyword})
	resp, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	res := []SearchResult{}

	doc.Find("table.sub-table tbody tr").Each(func(i int, s *goquery.Selection) {
		it := SearchResult{}

		it.Path, _ = url.QueryUnescape(s.Find("td a").AttrOr("href", ""))

		// path maybe urlencoded

		for i, n := range s.Find("td").EachIter() {
			switch i {
			case 0:
				it.Title = n.Text()
			case 2:
				text := n.Text()
				parts := strings.Split(text, " ")
				downloads, err := strconv.Atoi(parts[0])
				if err != nil {
					log.Fatalf("2nd td is %s", text)
				}
				it.Downloads = downloads
			}
		}
		res = append(res, it)
	})

	return res
}
