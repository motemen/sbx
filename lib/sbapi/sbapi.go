// Package sbapi provides unofficial Scrapbox API.
package sbapi

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"
)

type HTTPError http.Response

func (e HTTPError) Error() string {
	return fmt.Sprintf("%s: %s", e.Request.URL, e.Status)
}

type APIError struct {
	URL      string
	Response ErrorResponse
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Response.Name, e.Response.Message)
}

type ErrorResponse struct {
	Name    string      `json:"name"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

type PagesResponse struct {
	ProjectName string `json:"projectName"`
	Skip        int    `json:"skip"`
	Limit       int    `json:"limit"`
	Count       int    `json:"count"`
	Pages       []Page `json:"pages"`
}

type Project struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Image       string `json:"image,omitempty"`

	rawJSON json.RawMessage `json:"-"`
}

type Page struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Descriptions []string `json:"descriptions"`
	Accessed     Time     `json:"accessed"`
	Created      Time     `json:"created"`
	Updated      Time     `json:"updated"`

	rawJSON json.RawMessage `json:"-"`
}

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	var n int64
	err := json.Unmarshal(data, &n)
	if err != nil {
		return err
	}
	*t = Time(time.Unix(n, 0))
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Unix())
}

type options struct {
	sessionID string
	limit     uint
	origin    string
	host      string
	headers   map[string]string
}

type Option func(*options)

func WithSessionID(s string) Option {
	return func(o *options) {
		o.sessionID = s
	}
}

func WithLimit(n uint) Option {
	return func(o *options) {
		o.limit = n
	}
}

func WithOrigin(origin string) Option {
	return func(o *options) {
		o.origin = origin
	}
}

func WithHost(host string) Option {
	return func(o *options) {
		o.host = host
	}
}

func WithHeaders(h map[string]string) Option {
	return func(o *options) {
		o.headers = h
	}
}

func RequestJSON(path string, v interface{}, opts ...Option) error {
	wrapError := func(err error, message string) error {
		if err == nil {
			return err
		}

		if message == "" {
			return fmt.Errorf("%s %s: %w", "GET", path, err)
		}

		return fmt.Errorf("%s %s: %s: %w", "GET", path, message, err)
	}

	var opt options
	for _, o := range opts {
		o(&opt)
	}

	origin := "https://scrapbox.io"
	if opt.origin != "" {
		origin = opt.origin
	}

	url := origin + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range opt.headers {
		req.Header.Set(k, v)
	}

	if opt.host != "" {
		req.Host = opt.host
	}

	if opt.sessionID != "" {
		req.AddCookie(
			&http.Cookie{
				Name:  "connect.sid",
				Value: opt.sessionID,
			},
		)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	mediaType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))

	if resp.StatusCode >= 400 {
		if mediaType != "application/json" {
			return (*HTTPError)(resp)
		}

		var errResp ErrorResponse
		err = json.Unmarshal(data, &errResp)
		if err != nil {
			return wrapError(err, "json.Unmarshal")
		}

		err = APIError{
			URL:      url,
			Response: errResp,
		}
		return wrapError(err, "")
	}

	if strings.HasPrefix(mediaType, "text/") {
		switch x := v.(type) {
		case *string:
			*x = string(data)

		case *[]byte:
			*x = append([]byte{}, data...)

		case *interface{}:
			*x = append([]byte{}, data...)

		default:
			return fmt.Errorf("cannot handle %T for %s", v, mediaType)
		}

		return nil
	}

	err = json.Unmarshal(data, v)
	return wrapError(err, "json.Unmarshal")
}

func GetProject(projectName string, opts ...Option) (*Project, error) {
	var p Project

	err := RequestJSON(
		fmt.Sprintf("/api/projects/%s", projectName),
		&p,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func ListPages(projectName string, opts ...Option) ([]Page, error) {
	var opt options
	for _, o := range opts {
		o(&opt)
	}

	pages := []Page{}

	for {
		url := fmt.Sprintf("/api/pages/%s?skip=%d", projectName, len(pages))
		if opt.limit != 0 {
			url += fmt.Sprintf("&limit=%d", opt.limit)
		}

		var pagesResp PagesResponse
		err := RequestJSON(url, &pagesResp, opts...)
		if err != nil {
			return nil, err
		}

		pages = append(pages, pagesResp.Pages...)

		if (opt.limit != 0 && len(pages) >= int(opt.limit)) || len(pages) >= pagesResp.Count || len(pagesResp.Pages) == 0 {
			break
		}
	}

	return pages, nil
}
