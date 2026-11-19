package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Searcher interface {
	Search(ctx context.Context, query string, limit int) ([]Source, error)
}

type SearchClient struct {
	Endpoint string
	Key      string
	HTTP     *http.Client
}

func NewSearchClient(endpoint, key string) *SearchClient {
	return &SearchClient{Endpoint: endpoint, Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}
}

func (s *SearchClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	if k := strings.TrimSpace(s.Key); k != "" && k != "key" {
		req.Header.Set("Authorization", "Bearer "+k)
	}
	return s.HTTP.Do(req)
}

func (s *SearchClient) Search(ctx context.Context, query string, limit int) ([]Source, error) {
	if s.Endpoint == "" {
		return nil, fmt.Errorf("search engine belum dikonfigurasi")
	}
	u, err := url.Parse(s.Endpoint)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("query", query)
	q.Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		payload, _ := json.Marshal(map[string]any{"q": query, "query": query, "limit": limit})
		preq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.Endpoint, bytes.NewReader(payload))
		if err != nil {
			return nil, err
		}
		preq.Header.Set("Content-Type", "application/json")
		resp, err = s.do(preq)
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("search engine mengembalikan status %d", resp.StatusCode)
	}
	var body any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return ExtractSources(body, limit), nil
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}

func ExtractSources(v any, limit int) []Source {
	out := []Source{}
	seen := map[string]bool{}
	var walk func(n any)
	walk = func(n any) {
		if len(out) >= limit {
			return
		}
		switch t := n.(type) {
		case map[string]any:
			if link := firstString(t, "url", "link", "href", "source_url"); strings.HasPrefix(link, "http") {
				if !seen[link] {
					seen[link] = true
					title := firstString(t, "title", "name", "headline")
					if title == "" {
						title = link
					}
					snippet := firstString(t, "snippet", "content", "description", "body", "text", "summary")
					out = append(out, Source{Title: title, URL: link, Snippet: truncate(snippet, 400)})
				}
				return
			}
			keys := make([]string, 0, len(t))
			for k := range t {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				walk(t[k])
			}
		case []any:
			for _, x := range t {
				walk(x)
			}
		}
	}
	walk(v)
	return out
}
