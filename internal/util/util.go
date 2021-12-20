package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e2u/gp-scraper/internal/vars"
	gp_scraper "github.com/e2u/gp-scraper/vars"
	"github.com/k3a/html2text"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
)

func doRequest(ctx context.Context, method string, url string, headers map[string]string, body io.Reader) (int, []byte, error) {

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return 0, nil, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", vars.DefaultUserAgent)
	}

	if gp_scraper.Debug {
		db, _ := httputil.DumpRequest(req, true)
		fmt.Fprintf(os.Stdout, string(db))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, b, nil
}

func HTMLToText(html string) string {
	html2text.SetUnixLbr(true)
	return html2text.HTML2Text(html)
}

func HttpGet(ctx context.Context, url string, headers map[string]string) (int, []byte, error) {
	return doRequest(ctx, http.MethodGet, url, headers, nil)
}

func HttpPost(ctx context.Context, url string, headers map[string]string, body io.Reader) (int, []byte, error) {
	return doRequest(ctx, http.MethodPost, url, headers, body)
}

func IdentJSONString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func JSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
