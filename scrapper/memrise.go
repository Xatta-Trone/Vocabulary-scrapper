package scrapper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/xatta-trone/words-scrapper/model"
)

func ScrapMemrise(url string, options *model.Options) ([]model.Word, string, error) {

	words := []model.Word{}
	fileName := "default"
	var err error = nil
	var wordId int = 0

	c := colly.NewCollector(
		colly.AllowedDomains("app.memrise.com", "app.memrise.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
	)

	// Create another collector to scrape course details
	detailCollector := c.Clone()

	// Find the element with class word-list
	c.OnHTML("a.level.clearfix", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		// Print link
		// fmt.Printf("Link found: %s %s\n", link, e.Request.AbsoluteURL(link))
		detailCollector.Visit(e.Request.AbsoluteURL(link))

	})

	c.OnHTML("h1.course-name.sel-course-name", func(h *colly.HTMLElement) {
		title := strings.TrimSpace(h.Text)

		if len(title) > 0 {
			fileName = title
		}
	})

	// check error

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("There was an error, ", e)
		err = e
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Before making a request print "Visiting detailCollector..."
	detailCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting detailCollector", r.URL.String())
	})

	detailCollector.OnError(func(r *colly.Response, e error) {
		fmt.Println("There was an error, ", e)
		err = e
	})

	// scrap words

	detailCollector.OnHTML(".things", func(h *colly.HTMLElement) {
		h.DOM.Find(".thing").Each(func(i int, s *goquery.Selection) {
			word := model.Word{
				Word: strings.TrimSpace(strings.ReplaceAll(s.Find(".col_a").Text(), "\n", " ")),
			}

			if !options.NO_DEFINITION {
				word.Definition = strings.TrimSpace(strings.ReplaceAll(s.Find(".col_b").Text(), "\n", " "))
			}

			if !options.NO_ID {
				word.ID = wordId + 1
			}

			words = append(words, word)
			wordId++
		})

	})

	// Start scraping on https://app.memrise.com/course/5672405/barrons-gre-333-high-frequency-word/
	c.Visit(url)

	return words, fileName, err

}
