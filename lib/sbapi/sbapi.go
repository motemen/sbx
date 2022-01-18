// Package sbapi provides unofficial Scrapbox API.
package sbapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

	RawMessage json.RawMessage `json:"-"`
}

type Page struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Descriptions []string `json:"descriptions"`
	Accessed     Time     `json:"accessed"`
	Created      Time     `json:"created"`
	Updated      Time     `json:"updated"`

	RawMessage json.RawMessage `json:"-"`
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

func RequestJSON(url string, v interface{}, opts ...Option) error {
	wrapError := func(err error, message string) error {
		if err == nil {
			return err
		}

		if message == "" {
			return fmt.Errorf("%s %s: %w", "GET", url, err)
		}

		return fmt.Errorf("%s %s: %s: %w", "GET", url, message, err)
	}

	var opt options
	for _, o := range opts {
		o(&opt)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
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

	if resp.StatusCode >= 400 {
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

	err = json.Unmarshal(data, v)
	return wrapError(err, "json.Unmarshal")
}

func GetProject(projectName string, opts ...Option) (*Project, error) {
	var p Project

	err := RequestJSON(
		fmt.Sprintf("https://scrapbox.io/api/projects/%s", projectName),
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
		url := fmt.Sprintf("https://scrapbox.io/api/pages/%s?skip=%d", projectName, len(pages))
		if opt.limit != 0 {
			url += fmt.Sprintf("&limit=%d", opt.limit)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode >= 400 {
			var errResp ErrorResponse
			err = json.Unmarshal(data, &errResp)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("%s: %s", errResp.Name, errResp.Message)
		}

		var pagesResp PagesResponse
		err = json.Unmarshal(data, &pagesResp)
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
