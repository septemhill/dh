package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type uploadURLRequest struct {
	URL       string `json:"url"`
	ExpiredAt string `json:"expireAt"`
}

type UploadURLRequest struct {
	URL       string    `json:"url"`
	ExpiredAt time.Time `json:"expireAt"`
}

type UploadURLResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"shortUrl"`
}

func DecodeUploadRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	m := make(map[string]string)
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		return nil, err
	}

	expire, err := time.Parse(time.RFC3339, m["expireAt"])
	if err != nil {
		return nil, err
	}

	return UploadURLRequest{URL: m["url"], ExpiredAt: expire}, nil
}

type AccessURLRequest struct {
	ID string
}

type AccessURLResponse struct {
	URL string `json:"-"`
}

func (rsp *AccessURLResponse) Headers() http.Header {
	header := make(http.Header)
	header.Set("Location", rsp.URL)
	return header
}

func (rsp *AccessURLResponse) StatusCode() int {
	if len(rsp.URL) == 0 {
		return http.StatusNotFound
	}
	return http.StatusSeeOther
}

func DecodeAccessURLRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	var request AccessURLRequest
	request.ID = strings.TrimLeft(req.RequestURI, "/")

	return request, nil
}

func EncodeAccessURLResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	rsp := response.(AccessURLResponse)

	if len(rsp.URL) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	w.Header().Set("Location", rsp.URL)
	w.WriteHeader(http.StatusSeeOther)
	return nil
}
