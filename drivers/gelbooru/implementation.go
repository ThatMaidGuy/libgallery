package gelbooru

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ThatMaidGuy/libgallery"
	"github.com/ThatMaidGuy/libgallery/drivers/internal"
	"github.com/hashicorp/go-retryablehttp"
)

type implementation struct {
	host   string
	name   string
	client *http.Client
}

func New(name string, host string) libgallery.Driver {
	client := retryablehttp.NewClient()
	client.Logger = &internal.NoLogger{}
	return &implementation{
		client: client.StandardClient(),
		name:   name,
		host:   host,
	}
}

func (i *implementation) Search(query string, page uint64, limit uint64) ([]libgallery.Post, int, error) {
	if page > 200000/limit {
		return []libgallery.Post{}, 0, nil
	}

	const reqbase = "https://%s/index.php?page=dapi&s=post&q=index&limit=%v&tags=%s&pid=%v"
	url := fmt.Sprintf(reqbase, i.host, limit, url.QueryEscape(query), page)

	var response searchResponse
	err := internal.GetXML(url, i.client, &response)
	if err != nil {
		return []libgallery.Post{}, 0, err
	}

	/* The rule34.xxx API only has a success
	   value if there was an error. */
	if response.Success != nil {
		return []libgallery.Post{}, 0, response.Error
	}

	var posts []libgallery.Post

	for _, v := range response.Posts {
		ptime, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", v.CreatedAt)
		if err != nil {
			return []libgallery.Post{}, 0, err
		}

		var source []string
		if v.Source != "" {
			source = append(source, strings.TrimSpace(v.Source))
		}

		score, err := strconv.ParseInt(v.Score, 10, 64)
		if err != nil {
			return []libgallery.Post{}, 0, err
		}

		posts = append(posts, libgallery.Post{
			URL:    fmt.Sprintf("https://%s/index.php?page=post&s=view&id=%v", i.host, v.ID),
			Tags:   strings.TrimSpace(v.Tags),
			Date:   ptime,
			Source: source,
			ID:     v.ID,
			NSFW:   true,
			Score:  score,
		})
	}

	postsCount, _ := strconv.Atoi(response.Count)

	return posts, postsCount, err
}

func (i *implementation) File(id string) (libgallery.Files, error) {
	const reqbase = "https://%s/index.php?page=dapi&s=post&q=index&id="

	var response searchResponse
	err := internal.GetXML(fmt.Sprintf(reqbase, i.host)+id, i.client, &response)
	if err != nil {
		return []io.ReadCloser{}, err
	}

	// Same deal as before.
	if response.Success != nil {
		return []io.ReadCloser{}, response.Error
	}

	rc, err := internal.GetReadCloser(response.Posts[0].FileURL, i.client)
	if err != nil {
		return []io.ReadCloser{}, err
	}

	return []io.ReadCloser{rc}, nil
}

func (i *implementation) Name() string {
	return i.name
}

// Comments are broken on the API, I tried asking on their Discord
// about it twice with no response.
func (i *implementation) Comments(id string) ([]libgallery.Comment, error) {
	return []libgallery.Comment{}, nil
}
