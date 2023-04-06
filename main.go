package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xatta-trone/words-scrapper/model"
	"github.com/xatta-trone/words-scrapper/scrapper"
)

const (
	OUTPUT_FOLDER = "output"
)

func main() {

	fmt.Println("====== Welcome to Words Scrapper =======")
	fmt.Println("vocabulary.com | quizlet.com | memrise.com can be scrapped")

	url := flag.String("url", "", "The url to parse the words from")
	o := flag.String("o", "default", "The export file name [default is the word set title]")
	noDef := flag.Bool("no-def", false, "You don't want the definition to be parsed [default false]")
	noID := flag.Bool("no-id", false, "You don't want the word ID in the list [default false]")
	onlyCSV := flag.Bool("only-csv", false, "Export to CSV only [default false]")
	onlyJSON := flag.Bool("only-json", false, "Export to JSON only [default false]")
	wordsOnly := flag.Bool("words-only", false, "Only Scrap the words; no definition, no id will be included [default false]")
	flag.Parse()

	// construct the options

	options := model.Options{
		URL:           *url,
		OUTPUT:        *o,
		NO_DEFINITION: *noDef,
		NO_ID:         *noID,
		ONLY_WORD:     *wordsOnly,
		ONLY_CSV:      *onlyCSV,
		ONLY_JSON:     *onlyJSON,
	}

	if *wordsOnly {
		options.NO_DEFINITION = true
		options.NO_ID = true
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/scrap", func(ctx *gin.Context) {

		url := ctx.Query("url")

		// check if url is valid
		if !IsUrl(url) {
			fmt.Println("Please enter a valid URL with the flag -url=<your-URL-here>")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Please provide a valid url",
			})
			return
		}

		words, _, err := selectScrapper(url, &options)

		fmt.Println(words)
		fmt.Println(err)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		ctx.JSON(http.StatusOK, words)

	})
	r.Run("localhost:8080")

	// file name to be exported
	// defaultFileName := options.OUTPUT

	// words, file, err := selectScrapper(*url, &options)

	// if options.OUTPUT != "default" {
	// 	defaultFileName = options.OUTPUT
	// } else {
	// 	defaultFileName = file
	// }

	// if err != nil {
	// 	log.Fatalln("Could not parse url", err)
	// }

	// if len(words) > 0 {

	// 	if options.ONLY_CSV {
	// 		dumpCSV(words, defaultFileName)
	// 		return
	// 	}

	// 	if options.ONLY_JSON {
	// 		dumpJson(words, defaultFileName)
	// 		return
	// 	}
	// 	dumpCSV(words, defaultFileName)
	// 	dumpJson(words, defaultFileName)

	// }

}

func selectScrapper(url string, options *model.Options) (model.ResponseModel, string, error) {
	if strings.Contains(url, "vocabulary.com") {

		return scrapper.ScrapVocabulary(url, options)

	} else if strings.Contains(url, "quizlet.com") {

		return scrapper.DecideQuizletScrapper(url, options)
		// return scrapper.ScrapQuizlet(url, options)

	} else if strings.Contains(url, "memrise.com") {

		return scrapper.ScrapMemrise(url, options)

	} else {

		log.Fatal("The given url do not match vocabulary.com | quizlet.com | memrise.com")

		return model.ResponseModel{}, "default", errors.New("vocabulary.com | quizlet.com | memrise.com are only allowed")
	}
}

func dumpJson(words []model.Word, fileName string) {

	data, err := json.MarshalIndent(words, "", "\t")

	fileName = OUTPUT_FOLDER + "/" + fileName + ".json"

	if err != nil {
		log.Fatalln("Could not indent data", err)
	}

	ioutil.WriteFile(fileName, data, 0644)

	fmt.Println("File successfully written to ", fileName)
}

func dumpCSV(words []model.Word, fileName string) {

	fileName = OUTPUT_FOLDER + "/" + fileName + ".csv"

	file, err := os.Create(fileName)

	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	var data [][]string
	for _, record := range words {
		row := record.BuildCSV()
		data = append(data, row)
	}
	w.WriteAll(data)

	fmt.Println("File successfully written to ", fileName)
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
