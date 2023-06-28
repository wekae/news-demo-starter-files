package news

import "net/http"

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
