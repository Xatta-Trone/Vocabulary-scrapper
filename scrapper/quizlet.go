package scrapper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/gocolly/colly"
	"github.com/xatta-trone/words-scrapper/model"
)

func DecideQuizletScrapper(url string, options *model.Options) ([]model.Word, string, error) {

	words := []model.Word{}
	// indexes := map[string]int{}
	fileName := "default"
	var err error = nil

	if strings.Contains(url, "folders") && strings.Contains(url, "sets") {
		urls, errs := GetUrlMaps(url)

		if errs != "" {
			fmt.Println(errs)
		}

		for key, val := range urls {
			fmt.Println(key, val)

			wds, _, _ := ScrapQuizlet(key, options,val)

			words = append(words, wds...)

		}
		return words, fileName, err

	} else {
		return ScrapQuizlet(url, options, 1)
	}

}

func GetUrlMaps(url string) (map[string]int, string) {
	indexes := map[string]int{}
	err := ""
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			// fmt.Println(string(r.Body))

			if r.StatusCode != http.StatusOK {
				err = r.Status
			}

			root := r.HTMLDoc.Find(".FolderPageSetsList-setsFeed")

			if root.Length() == 1 {
				// it will go through each group
				sets := root.Find(".UISetCard")
				length := sets.Length()

				// fmt.Println(length)

				if length > 0 {
					sets.Each(func(i int, s *goquery.Selection) {
						// get the url
						setUrl := s.Find(".UIBaseCardHeader a").AttrOr("href", "")
						fmt.Println(setUrl)

						indexes[setUrl] = length
						length--

					})
				}

			}
		},
	}).Start()

	return indexes, err
}

func ScrapQuizlet(url string, options *model.Options, groupId int) ([]model.Word, string, error) {

	words := []model.Word{}
	// indexes := map[string]int{}
	fileName := "default"
	var err error = nil
	// indexBeforeLogin := 0

	c := colly.NewCollector(
		colly.AllowedDomains("quizlet.com", "www.quizlet.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
	)

	// check if its a folder
	c.OnHTML(".FolderPageSetsList-setsFeed", func(h *colly.HTMLElement) {
		// it will go through each group
		sets := h.DOM.Find(".UIDiv")
		length := sets.Length()

		if length > 0 {
			sets.Each(func(i int, s *goquery.Selection) {
				fmt.Println(s.Find(".UIBaseCardHeader a").Attr("href"))
			})
		}

	})

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
				Group: groupId,
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
				word2.Group = groupId

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
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, ":", "")
		if len(title) > 0 {
			fileName = title
		}
	})

	// check error
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("There was an error, ", e.Error())
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
