package client

import (
	"bytes"
	"context"
	"devices-test/configs"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	Client   http.Client
	Context  context.Context
	question int
	Cookies  []*http.Cookie
	Headers  map[string]string
	Configs  *configs.Config
}

func NewClient(ctx context.Context, cfg *configs.Config) *Client {
	cookieJar, _ := cookiejar.New(nil)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return &Client{
		Client:   http.Client{Jar: cookieJar},
		Context:  ctx,
		Headers:  headers,
		question: 0,
		Configs:  cfg,
	}
}

func (c *Client) Run() error {
	req, _ := http.Get(c.Configs.BaseURL)
	for {
		c.question++
		data, err := c.parse(req)
		if err != nil {
			return err
		}
		body, err := c.sendRequest(data)
		if strings.Contains(body, "Test successfully passed") {
			return nil
		}
	}
}

func (c *Client) endPointURL() string {
	return c.Configs.BaseURL + c.Configs.EndPath + strconv.Itoa(c.question)
}

func (c *Client) sendRequest(data url.Values) (string, error) {
	req, err := http.NewRequest(http.MethodPost, c.endPointURL(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	for _, v := range c.Cookies {
		req.AddCookie(v)
	}
	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("not successfull request: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) parse(r *http.Response) (url.Values, error) {
	req, err := http.NewRequestWithContext(c.Context, http.MethodGet, c.endPointURL(), nil)
	if err != nil {
		log.Fatal("cannot read request: ", err)
	}
	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}
	c.Cookies = r.Cookies()
	for _, v := range c.Cookies {
		req.AddCookie(v)
	}

	var resp *http.Response
	resp, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("not successfull request: %w", err)
	}
	defer resp.Body.Close()

	m := map[string]string{}
	z := html.NewTokenizer(resp.Body)
	var selectName string
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			data := url.Values{}
			for k, v := range m {
				data.Add(k, v)
			}
			return data, err
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "input" {
				if t.Attr[0].Val == "text" {
					m[t.Attr[1].Val] = "test"
				}
				if t.Attr[0].Val == "radio" {
					newRadioName := t.Attr[1].Val
					if val, ok := m[newRadioName]; ok {
						if len(t.Attr[2].Val) > len(val) {
							m[newRadioName] = t.Attr[2].Val
						}
					} else {
						m[newRadioName] = t.Attr[2].Val
					}
				}
			}
			if t.Data == "select" {
				m[t.Attr[0].Val] = ""
				selectName = t.Attr[0].Val
			}
			if t.Data == "option" {
				if val, ok := m[selectName]; ok {
					if len(t.Attr[0].Val) > len(val) {
						m[selectName] = t.Attr[0].Val
					}
				}
			}
		}
	}
}
