package nix

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	str "github.com/vestlog/nix/pkg/structs"
)

type APIClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

func (c *APIClient) Get(url string) ([]byte, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalPosts(data []byte) ([]str.Post, error) {
	posts := make([]str.Post, 0)
	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
	}
	return posts, nil
}

func unmarshalComments(data []byte) ([]str.Comment, error) {
	comments := make([]str.Comment, 0)
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
	}
	return comments, nil
}

func (c *APIClient) GetPosts(userID int) ([]str.Post, error) {
	req := fmt.Sprintf("posts?userId=%d", userID)
	data, err := c.Get(req)
	if err != nil {
		return nil, fmt.Errorf("error getting url: %w", err)
	}
	posts, err := unmarshalPosts(data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}
	return posts, nil
}

func (c *APIClient) GetComments(postID int) ([]str.Comment, error) {
	req := fmt.Sprintf("comments?postId=%d", postID)
	data, err := c.Get(req)
	if err != nil {
		return nil, fmt.Errorf("error getting url: %w", err)
	}
	comments, err := unmarshalComments(data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}
	return comments, nil
}

func CreateAPIClient(dsn string, url string) (*APIClient, error) {
	return &APIClient{
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
		BaseURL: url,
	}, nil
}
