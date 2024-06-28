package helpers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "genie"}

func ScrapeURL(startURL, element, pagination string, limit int) (string, error) {
	var allData []string
	currentPage := startURL
	pagesScraped := 0

	parsedBaseURL, err := url.Parse(startURL)
	if err != nil {
		return "", err
	}

	for {
		if limit > 0 && pagesScraped >= limit {
			break
		}

		doc, err := fetchDocument(currentPage)
		if err != nil {
			return "", err
		}

		var data []string
		if element != "" {
			data = extractElements(doc, element)
		} else {
			data = extractAllContent(doc)
		}
		allData = append(allData, data...)
		pagesScraped++

		if pagination == "" {
			break
		}

		nextPage := getNextPage(doc, pagination, parsedBaseURL)
		if nextPage == "" {
			break
		}

		currentPage = nextPage
	}

	return strings.Join(allData, "\n"), nil
}

func fetchDocument(pageURL string) (*goquery.Document, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func extractElements(doc *goquery.Document, element string) []string {
	var data []string
	doc.Find(element).Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if element == "a" {
			href, exists := s.Attr("href")
			if exists {
				text = fmt.Sprintf("%s (%s)", text, href)
			}
		}
		data = append(data, text)
	})
	return data
}

func extractAllContent(doc *goquery.Document) []string {
	return []string{doc.Text()}
}

func getNextPage(doc *goquery.Document, pagination string, baseURL *url.URL) string {
	nextPage := ""
	doc.Find(pagination).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			nextPage = href
			return
		}
	})

	if nextPage != "" {
		nextPageURL, err := baseURL.Parse(nextPage)
		if err == nil {
			return nextPageURL.String()
		}
	}

	return ""
}
