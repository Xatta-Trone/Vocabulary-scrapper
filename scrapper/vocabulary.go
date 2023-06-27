package scrapper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/xatta-trone/words-scrapper/model"
)

func ScrapVocabulary(url string, options *model.Options) (model.ResponseModel, string, error) {

	// words := []model.Word{}
	fileName := "default"
	var err error = nil

	var finalResult model.ResponseModel

	finalResult.FolderURL = url

	// since each vocabulary url has one list
	var singleResponse model.SingleResponseModel
	singleResponse.GroupId = 1
	singleResponse.URL = url

	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			// fmt.Println(string(r.Body))

			if r.StatusCode != http.StatusOK {
				fmt.Println("There was an error, ", r.Status)
				err = fmt.Errorf("%s", r.Status)
			}

			// get the title 

			title := r.HTMLDoc.Find("h1.title").Text()
			title = strings.TrimSpace(title)

			fmt.Print(title)

			if len(title) > 0 {
				fileName = title
				singleResponse.Title = title
			}



			// get the words


			root := r.HTMLDoc.Find("ol.wordlist")

			fmt.Println(root.Length())

			if root.Length() == 1 {

				root.Children().Each(func(i int, s *goquery.Selection) {
					wordCheck := s.AttrOr("word", "")

					fmt.Println(wordCheck)

					if wordCheck != "" {
						word := strings.TrimSpace(strings.ReplaceAll(wordCheck, "\n", " "))
						singleResponse.Words = append(singleResponse.Words, word)

					}

				})

			}

		},
	}).Start()

	// c := colly.NewCollector(
	// 	colly.AllowedDomains("www.vocabulary.com"),
	// 	colly.UserAgent("Mozilla/5.0 (X11; Linux i686; rv:109.0) Gecko/20100101 Firefox/114.0"),
	// )

	// // Find the element with class word-list
	// c.OnHTML("ol.wordlist", func(e *colly.HTMLElement) {

	// 	e.DOM.Children().Each(func(i int, s *goquery.Selection) {
	// 		// we are inside each list element
	// 		// <li class="entry learnable" id="entry1" word="estranged" freq="2906.44" lang="en">
	// 		// <a class="word" href="/dictionary/estranged" title="caused to be unloved"><span class="count"></span> estranged</a>
	// 		// <div class="definition" title="This word is learnable">caused to be unloved</div>
	// 		// </li>

	// 		// check if word is not null or exists
	// 		wordCheck := s.AttrOr("word", "")

	// 		if wordCheck != "" {
	// 			singleResponse.Words = append(singleResponse.Words, strings.TrimSpace(strings.ReplaceAll(s.AttrOr("word", ""), "\n", " ")))

	// 			// word := model.Word{
	// 			// 	Word: strings.TrimSpace(strings.ReplaceAll(s.AttrOr("word", ""), "\n", " ")),
	// 			// }

	// 			// if !options.NO_DEFINITION {
	// 			// 	word.Definition = strings.TrimSpace(strings.ReplaceAll(s.Find(".definition").Text(), "\n", " "))
	// 			// }

	// 			// if !options.NO_ID {
	// 			// 	word.ID = i + 1
	// 			// }

	// 			// words = append(words, word)

	// 		}

	// 	})

	// })

	// c.OnHTML("h1.title", func(h *colly.HTMLElement) {
	// 	title := strings.TrimSpace(h.Text)

	// 	if len(title) > 0 {
	// 		fileName = title
	// 		singleResponse.Title = title
	// 	}
	// })

	// // check error

	// c.OnError(func(r *colly.Response, e error) {
	// 	fmt.Println("There was an error, ", e)
	// 	err = e
	// })

	// // Before making a request print "Visiting ..."
	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL.String())
	// })

	// // Start scraping on https://www.vocabulary.com/lists/7200740
	// c.Visit(url)

	finalResult.Sets = append(finalResult.Sets, singleResponse)

	return finalResult, fileName, err

}
