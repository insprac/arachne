package crawler

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

type ProcessFunc func(doc Document) error

type Document struct {
	Uri   string `json:"uri"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var ignoredTags = []string{"script", "style", "header", "footer", "nav", "aside"}

func Crawl(startingURLs, allowedURLPrefixes []string, process ProcessFunc) error {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	scrapedPages := make(map[string]bool)

	var title string

	c.OnHTML("title", func(e *colly.HTMLElement) {
		title = e.Text
	})

	c.OnResponse(func(r *colly.Response) {
		uri := r.Request.URL.String()

		queryDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
		if err != nil {
			log.Fatal("Failed to parse response body:", err)
			return
		}

		// Strip ignored tags to avoid repeated info such as navigation and footer
		for _, tag := range ignoredTags {
			queryDoc.Find(tag).Each(func(i int, s *goquery.Selection) {
				s.Remove()
			})
		}

		manipulatedBody, err := queryDoc.Html()
		if err != nil {
			log.Fatal("Failed to generate manipulated body:", err)
			return
		}

		body := convertHTMLToMarkdown(string(manipulatedBody))

		doc := Document{
			Uri:   uri,
			Title: title,
			Body:  body,
		}

		err = process(doc)
		if err != nil {
			fmt.Printf("Error processing crawled document: %v", err)
		}
	})

	// Visit links on the page
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		u, err := url.Parse(link)
		if err != nil {
			return
		}

		absoluteLink := e.Request.AbsoluteURL(u.String())

		if isUrlAllowed(absoluteLink, allowedURLPrefixes) && !scrapedPages[absoluteLink] {
			scrapedPages[absoluteLink] = true

			err = c.Visit(absoluteLink)
			if err != nil {
				return
			}
		}
	})

	for _, startingURL := range startingURLs {
		err := c.Visit(startingURL)
		if err != nil {
			return err
		}
	}

	c.Wait()
	return nil
}

func isUrlAllowed(link string, allowedURLPrefixes []string) bool {
	for _, prefix := range allowedURLPrefixes {
		if strings.HasPrefix(link, prefix) {
			return true
		}
	}
	return false
}

func convertHTMLToMarkdown(html string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(html)
	if err != nil {
		return ""
	}

	return markdown
}
