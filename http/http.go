package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// GetPageCtx gets the body at the given url and returns it as an io.Reader
// using context.
func GetPageCtx(ctx context.Context, url string) (*bytes.Buffer, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error %d %s", res.StatusCode, res.Status)
	}
	var b bytes.Buffer
	_, err = io.Copy(&b, res.Body)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
