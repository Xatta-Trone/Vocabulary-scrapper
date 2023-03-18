package scrapper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/xatta-trone/words-scrapper/model"
)

func ScrapQuizlet(url string, options *model.Options) ([]model.Word, string, error) {

	words := []model.Word{}
	fileName := "default"
	var err error = nil
	// indexBeforeLogin := 0

	c := colly.NewCollector(
		colly.AllowedDomains("quizlet.com", "www.quizlet.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
	)

	// Find the element with class SetPageTerms-term

	c.OnHTML(".SetPageTerms-termsWrapper", func(e *colly.HTMLElement) {
		currentId := 0
		// find total element in this set
		// totalSet := strings.TrimSpace(e.DOM.Children().Find(".t27kl0s").Text())
		// fmt.Println("total words,", totalSet)

		// find the free words
		e.DOM.Children().Find(".SetPageTerms-term").Each(func(i int, s *goquery.Selection) {

			word := model.Word{
				Word: s.Children().Find(".SetPageTerm-wordText").Text(),
			}

			if !options.NO_DEFINITION {
				word.Definition = s.Children().Find(".SetPageTerm-definitionText").Text()
			}

			if !options.NO_ID {
				word.ID = currentId + 1
			}

			// fmt.Println(word)

			words = append(words, word)
			currentId++

		})

		// now go for remaining words
		word2 := model.Word{}

		e.DOM.Find("div[style=\"display:none\"]").Children().Each(func(i int, s *goquery.Selection) {

			// set the word // word is in every even number element
			str := strings.TrimSpace(strings.ReplaceAll(s.Text(), "\n", " "))
			if i == 0 || i%2 == 0 {
				word2.Word = str

				if !options.NO_ID {
					word2.ID = currentId + 1
				}

				currentId++
			}

			// set the definition // definition is in every odd number of element
			if i%2 == 1 {
				if !options.NO_DEFINITION {
					word2.Definition = str
				}
				// fmt.Println(word2)
				words = append(words, word2)
				// set the model empty for the next word
				word2 = model.Word{}
			}

		})

	})

	c.OnHTML("div.SetPage-titleWrapper", func(h *colly.HTMLElement) {
		title := strings.TrimSpace(h.Text)
		title = strings.ReplaceAll(title," ","-")
		title = strings.ReplaceAll(title,":","")
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

	// Start scraping on https://quizlet.com/130371046/gre-flash-cards/
	c.Visit(url)

	return words, fileName, err

}
