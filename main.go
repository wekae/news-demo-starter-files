package main

import (
	"fmt"
	"github.com/freshman-tech/news-demo-starter-files/news"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/freshman-tech/news-demo-starter-files/news"
	_ "time"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Search Query is: ", searchQuery)
	fmt.Println("Page is: ", page)
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

	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
