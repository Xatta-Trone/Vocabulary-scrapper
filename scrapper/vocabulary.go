package scrapper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/xatta-trone/words-scrapper/model"
)

func ScrapVocabulary(url string, options *model.Options) ([]model.Word, string, error) {

	words := []model.Word{}
	fileName := "default"
	var err error = nil

	c := colly.NewCollector(
		colly.AllowedDomains("www.vocabulary.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
	)

	// Find the element with class word-list
	c.OnHTML("ol.wordlist", func(e *colly.HTMLElement) {

		e.DOM.Children().Each(func(i int, s *goquery.Selection) {
			// we are inside each list element
			// <li class="entry learnable" id="entry1" word="estranged" freq="2906.44" lang="en">
			// <a class="word" href="/dictionary/estranged" title="caused to be unloved"><span class="count"></span> estranged</a>
			// <div class="definition" title="This word is learnable">caused to be unloved</div>
			// </li>

			// check if word is not null or exists
			wordCheck := s.AttrOr("word", "")

			if wordCheck != "" {
				// word := model.Word{
				// 	ID:         i + 1,
				// 	Word:       s.AttrOr("word", ""),
				// 	Definition: s.Find(".definition").Text(),
				// }

				word := model.Word{
					Word: strings.TrimSpace(strings.ReplaceAll(s.AttrOr("word", ""), "\n", " ")),
				}

				if !options.NO_DEFINITION {
					word.Definition = strings.TrimSpace(strings.ReplaceAll(s.Find(".definition").Text(), "\n", " "))
				}

				if !options.NO_ID {
					word.ID = i + 1
				}

				words = append(words, word)

			}

		})

	})

	c.OnHTML("h1.title", func(h *colly.HTMLElement) {
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

	// Start scraping on https://www.vocabulary.com/lists/7200740
	c.Visit(url)

	return words, fileName, err

}
