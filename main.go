package main

import (
	"bytes"
	"fmt"
	"github.com/freshman-tech/news-demo-starter-files/news"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	_ "strconv"
	"time"

	_ "github.com/freshman-tech/news-demo-starter-files/news"
)

type Search struct {
	Query      string
	NextPage   int
	TotalPages int
	Results    *news.Results
}

var tpl = template.Must(template.ParseFiles("index.html"))

/*
refactored indexHandler so that the template is no longer executed directly to ResponseWriter
*/
func indexHandler(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
}

/*
*
The searchHandler function has been changed.
It now accepts a pointer to news.Client and returns an anonymous function which satisfies the http.HandlerFunc type.
This function closes over the newsapi parameter which means it will have access to it whenever it is called.
*/
func searchHandler(newsApi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" {
			page = "1"
		}

		results, err := newsApi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nextPage, err := strconv.Atoi(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		search := &Search{
			Query:      searchQuery,
			NextPage:   nextPage,
			TotalPages: int(math.Ceil(float64(results.TotalResults) / float64(newsApi.PageSize))),
			Results:    results,
		}

		/*
			The template is first executed into an empty buffer so that we can check for errors.
			After that, the buffer is written to the ResponseWriter.
			If we execute the template directly on ResponseWriter,
			we won’t be able to check for errors so this is a better way to do it.
		*/
		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, search)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)

		fmt.Printf("%+v", results)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatalln("ENV: NEWS_API_KEY must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsApi := news.NewClient(myClient, apiKey, 20)

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/search", searchHandler(newsApi))
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
