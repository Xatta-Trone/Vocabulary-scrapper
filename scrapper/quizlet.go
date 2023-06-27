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

func DecideQuizletScrapper(url string, options *model.Options) (model.ResponseModel, string, error) {

	// words := []model.Word{}
	// indexes := map[string]int{}
	fileName := "default"
	var err error = nil

	var finalResult model.ResponseModel
	finalResult.FolderURL = url

	if strings.Contains(url, "folders") && strings.Contains(url, "sets") {
		urls, folderName, errs := GetUrlMaps(url)

		if errs != "" {
			fmt.Println(errs)
		}

		fmt.Println(folderName)

		for _, set := range urls {
			fmt.Println(set.ID, set.Url)

			wds, file, errs := ScrapQuizlet(set.Url, options, set.ID)

			finalResult.Sets = append(finalResult.Sets, wds)
			fileName = file
			err = errs

		}
		return finalResult, fileName, err

	} else {
		data, fileN, errs := ScrapQuizlet(url, options, 1)

		fileName = fileN
		err = errs
		finalResult.Sets = append(finalResult.Sets, data)

		return finalResult, fileName, err
	}

}

func GetUrlMaps(url string) ([]model.QuizletFolder, string, string) {
	indexes := []model.QuizletFolder{}
	err := ""
	folderName := ""
	geziyor.NewGeziyor(&geziyor.Options{
		StartRequestsFunc: func(g *geziyor.Geziyor) {
			g.GetRendered(url, g.Opt.ParseFunc)
		},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			// fmt.Println(string(r.Body))

			if r.StatusCode != http.StatusOK {
				err = r.Status
			}

			// find the title
			titleText := r.HTMLDoc.Find("div.DashboardHeaderTitle-main").Text()
			title := strings.TrimSpace(titleText)
			// title = strings.ReplaceAll(title, " ", "-")
			title = strings.ReplaceAll(title, ":", "")
			if len(title) > 0 {
				folderName = title
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
						// temp data
						model := model.QuizletFolder{ID: length, Url: setUrl}

						indexes = append(indexes, model)
						length--

					})
				}

			}
		},
	}).Start()

	return indexes, folderName, err
}

func ScrapQuizlet(url string, options *model.Options, groupId int) (model.SingleResponseModel, string, error) {

	// words := []model.Word{}
	// indexes := map[string]int{}
	fileName := "default"
	var err error = nil
	// indexBeforeLogin := 0
	var singleResponse model.SingleResponseModel
	singleResponse.GroupId = groupId
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

			root := r.HTMLDoc.Find(".SetPage-setContentWrapper")
			// for {
			// 	// chromedp.WaitVisible("button[aria-label='See more']")
			// 	s := root.Find("button").Last()
			// 	fmt.Println(s.Text())
			// 	if s.Text() == "See More" {
			// 		chromedp.Click(`*/deep/[value="See More"]`, chromedp.BySearch)
			// 		time.Sleep(time.Duration(time.Second * 5))
			// 		continue
			// 	} else {
			// 		break
			// 	}

			// }

			if root.Length() == 1 {
				// it will go through each group
				rootSet := root.Find(".SetPageTerms-termsWrapper")
				sets := rootSet.Find(".SetPageTerms-term")
				length := sets.Length()

				fmt.Println(length)

				if length > 0 {
					sets.Each(func(i int, s *goquery.Selection) {
						word := s.Children().Find(".SetPageTerm-wordText").Text()
						singleResponse.Words = append(singleResponse.Words, word)
					})
				}

				// hidden sets
				hidden := root.Find("div[style=\"display:none\"]").Children()

				fmt.Println("hidden set", hidden.Length())

				hidden.Each(func(i int, s *goquery.Selection) {
					// set the word // word is in every even number element
					str := strings.TrimSpace(strings.ReplaceAll(s.Text(), "\n", " "))
					if i == 0 || i%2 == 0 {
						singleResponse.Words = append(singleResponse.Words, str)

						// word2.Word = str
						// word2.Group = groupId

						// if !options.NO_ID {
						// 	word2.ID = currentId + 1
						// }

						// currentId++
					}

				})

				// find the title
				titleText := root.Find("div.SetPage-breadcrumbTitleWrapper").Text()
				title := strings.TrimSpace(titleText)
				// title = strings.ReplaceAll(title, " ", "-")
				title = strings.ReplaceAll(title, ":", "")
				if len(title) > 0 {
					fileName = title
					singleResponse.Title = title
				}

			}

		},
	}).Start()

	return singleResponse, fileName, err

}

// func ScrapQuizlet(url string, options *model.Options, groupId int) (model.SingleResponseModel, string, error) {

// 	// words := []model.Word{}
// 	// indexes := map[string]int{}
// 	fileName := "default"
// 	var err error = nil
// 	// indexBeforeLogin := 0
// 	var singleResponse model.SingleResponseModel
// 	singleResponse.GroupId = groupId
// 	singleResponse.URL = url

// 	c := colly.NewCollector(
// 		colly.AllowedDomains("quizlet.com", "www.quizlet.com"),
// 		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0"),
// 	)

// 	// check if its a folder
// 	c.OnHTML(".FolderPageSetsList-setsFeed", func(h *colly.HTMLElement) {
// 		// it will go through each group
// 		sets := h.DOM.Find(".UIDiv")
// 		length := sets.Length()

// 		if length > 0 {
// 			sets.Each(func(i int, s *goquery.Selection) {
// 				fmt.Println(s.Find(".UIBaseCardHeader a").Attr("href"))
// 			})
// 		}

// 	})

// 	// Find the element with class SetPageTerms-term
// 	c.OnHTML(".SetPageTerms-termsWrapper", func(e *colly.HTMLElement) {
// 		// currentId := 0
// 		// find total element in this set
// 		// totalSet := strings.TrimSpace(e.DOM.Children().Find(".t27kl0s").Text())
// 		// fmt.Println("total words,", totalSet)

// 		// find the free words
// 		e.DOM.Children().Find(".SetPageTerms-term").Each(func(i int, s *goquery.Selection) {

// 			singleResponse.Words = append(singleResponse.Words, s.Children().Find(".SetPageTerm-wordText").Text())

// 			// word := model.Word{
// 			// 	Word: s.Children().Find(".SetPageTerm-wordText").Text(),
// 			// 	Group: groupId,
// 			// }

// 			// if !options.NO_DEFINITION {
// 			// 	word.Definition = s.Children().Find(".SetPageTerm-definitionText").Text()
// 			// }

// 			// if !options.NO_ID {
// 			// 	word.ID = currentId + 1
// 			// }

// 			// // fmt.Println(word)

// 			// words = append(words, word)
// 			// currentId++

// 		})

// 		// now go for remaining words
// 		// word2 := model.Word{}

// 		e.DOM.Find("div[style=\"display:none\"]").Children().Each(func(i int, s *goquery.Selection) {

// 			// set the word // word is in every even number element
// 			str := strings.TrimSpace(strings.ReplaceAll(s.Text(), "\n", " "))
// 			if i == 0 || i%2 == 0 {
// 				singleResponse.Words = append(singleResponse.Words, str)

// 				// word2.Word = str
// 				// word2.Group = groupId

// 				// if !options.NO_ID {
// 				// 	word2.ID = currentId + 1
// 				// }

// 				// currentId++
// 			}

// 			// set the definition // definition is in every odd number of element
// 			// if i%2 == 1 {
// 			// 	if !options.NO_DEFINITION {
// 			// 		word2.Definition = str
// 			// 	}
// 			// 	// fmt.Println(word2)
// 			// 	words = append(words, word2)
// 			// 	// set the model empty for the next word
// 			// 	word2 = model.Word{}
// 			// }

// 		})

// 	})

// 	c.OnHTML("div.SetPage-titleWrapper", func(h *colly.HTMLElement) {
// 		title := strings.TrimSpace(h.Text)
// 		title = strings.ReplaceAll(title, " ", "-")
// 		title = strings.ReplaceAll(title, ":", "")
// 		if len(title) > 0 {
// 			fileName = title
// 			singleResponse.Title = title
// 		}
// 	})

// 	// check error
// 	c.OnError(func(r *colly.Response, e error) {
// 		fmt.Println("There was an error, ", e.Error())
// 		err = e
// 	})

// 	// Before making a request print "Visiting ..."
// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting", r.URL.String())
// 	})

// 	// Start scraping on https://quizlet.com/130371046/gre-flash-cards/
// 	c.Visit(url)

// 	return singleResponse, fileName, err

// }
