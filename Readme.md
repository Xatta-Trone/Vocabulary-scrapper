# Vocabulary scrapper
## Scrap words from | vocabulary.com | quizlet.com | memrise.com with ease. 

I built it as a hobby project while preparing for GRE

---

```bash
# Basic example
go run main.go -url https://app.memrise.com/course/5672405/barrons-gre-333-high-frequency-word/
go run main.go -url https://quizlet.com/130371046/gre-flash-cards/
go run main.go -url https://www.vocabulary.com/lists/7200740

# Additional flags
go run main.go -url https://www.vocabulary.com/lists/7200740 -o myCustomFileName -words-only


```

## Available options 

```
====== Welcome to Words Scrapper =======
vocabulary.com | quizlet.com | memrise.com can be scrapped
Usage:   
  -no-def
        You don't want the definition to be parsed [default false]
  -no-id
        You don't want the word ID in the list [default false]
  -o string
        The export file name [default is the word set title] (default "default")   
  -only-csv
        Export to CSV only [default false]
  -only-json
        Export to JSON only [default false]
  -url string
        The url to parse the words from
  -words-only
        Only Scrap the words; no definition, no id will be included [default false]
```


**The output file will be saved in the `output` folder**





## Packages used 
1. Colly: https://github.com/gocolly/colly
2. GoQuery: https://github.com/PuerkitoBio/goquery
