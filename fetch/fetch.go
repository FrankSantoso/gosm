package fetch

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseurl    = "http://nominatim.openstreetmap.org/search.php"
	deftimeout = 5 * time.Second
)

type SearchQuery struct {
	string
}

// buildQuery builds valid url query from search query string
func (s *SearchQuery) buildQuery() string {
	var query strings.Builder
	components := []string{baseurl, "?q=", url.QueryEscape(s.string)}
	for _, v := range components {
		query.Write([]byte(v))
	}
	return query.String()
}

func Fetch(ctx context.Context, query string) (io.ReadCloser, error) {
	newctx, _ := context.WithTimeout(ctx, deftimeout)
	sq := &SearchQuery{query}
	req, err := http.NewRequestWithContext(newctx, "GET", sq.buildQuery(), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	return resp.Body, nil
}
