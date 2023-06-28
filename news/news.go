package news

import (
	"net/http"
	"time"
)

/**
The Client struct represents the client for working with the News API.
The httpClient field points to the HTTP client that should be used to make requests,
apiKey field holds the API key while the PageSize field holds the number of results to
return per page (maximum of 100).
The NewClient() function creates and returns a new Client instance for making requests to the News API.
*/

type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}

type Article struct {
	Source struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}
