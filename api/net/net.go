package net

import (
	"code.google.com/p/go.net/publicsuffix"
	"github.com/joffrey-bion/gorol/api/parser"
	"github.com/joffrey-bion/gorol/log"
	"github.com/joffrey-bion/gorol/model"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	fakeUserAgent string = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.143 Safari/537.36"
)

var (
	state  model.AccountState
	client *http.Client
)

func init() {
	// set http client
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		panic(err)
	}
	client = &http.Client{Jar: jar}
}

func pageUrl(base string, page string, query url.Values) string {
	url := base + "?p=" + page
	if len(query) > 0 {
		url += "&" + query.Encode()
	}
	return url
}

func Get(base, page string, query url.Values) (string, error) {
	url := pageUrl(base, page, query)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	return execute(request)
}

func Post(base, page string, form url.Values) (string, error) {
	url := pageUrl(base, page, nil)
	request, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return execute(request)
}

func execute(request *http.Request) (string, error) {
	request.Header.Set("User-Agent", fakeUserAgent)
	log.Df("url=%s", request.URL)
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	respBody, err := parser.String(resp)
	if err != nil {
		return "", err
	}
	return respBody, nil
}
