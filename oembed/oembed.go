package oembed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// replacements to convert patterns to regexes
var (
	su2re1 = regexp.MustCompile("^(https?://[^/]*?)\\*(.+)$")
	su2re2 = regexp.MustCompile("^(https?://[^/]*?/.*?)\\*(.+)$")
	su2re3 = regexp.MustCompile("^(https?://.*?)\\*$")
	su2re4 = regexp.MustCompile("^http://")
)

// Oembed contains list of available oembed items (official endpoints)
type Oembed struct {
	items []*Item
}

// Endpoint contains single endpoint to check against
type Endpoint struct {
	URL       string   `json:"url"`
	Discovery bool     `json:"discovery,omitempty"`
	Schemes   []string `json:"schemes,omitempty"`
}

// Provider contains a single provider which can have multiple endpoints
type Provider struct {
	Name      string     `json:"provider_name"`
	URL       string     `json:"provider_url"`
	Endpoints []Endpoint `json:"endpoints"`
}

// Item contains data for a schema
type Item struct {
	IsEndpointURLComplete bool
	EndpointURL           string
	ProviderName          string
	ProviderURL           string
	regex                 *regexp.Regexp
}

// Options contains options to pass to FetchOembed method
type Options struct {
	Client         *http.Client
	URL            string
	MaxWidth       int
	MaxHeight      int
	AcceptLanguage string
	ExtraOpts      map[string]string
}

// ComposeURL returns url of oembed resource ready to be queried
func (item *Item) ComposeURL(u string) string {
	if item.IsEndpointURLComplete {
		return item.EndpointURL
	}

	return item.EndpointURL + url.QueryEscape(u)
}

func (item *Item) parseOembed(u string, resp *http.Response) (*Info, error) {
	var err error
	if resp.StatusCode > 200 {
		return &Info{Status: resp.StatusCode}, nil
	}

	reader := io.LimitReader(resp.Body, 40000) // 40 KB max

	info := NewInfo()
	err = info.FillFromJSON(reader)

	if err != nil {
		return nil, err
	}

	if len(info.URL) == 0 {
		info.URL = u
	}

	if len(info.ProviderURL) == 0 {
		info.ProviderURL = item.ProviderURL
	}

	if len(info.ProviderName) == 0 {
		info.ProviderName = item.ProviderName
	}

	return info, nil
}

func buildQueryString(params map[string]string) string {
	var b bytes.Buffer

	for key, value := range params {
		b.WriteString(fmt.Sprintf("&%s=%s", key, value))
	}

	return b.String()
}

// FetchOembed return oembed info from an url containing it
func (item *Item) FetchOembed(opts Options) (*Info, error) {
	resURL := item.ComposeURL(opts.URL)

	if opts.MaxWidth > 0 {
		resURL = fmt.Sprintf("%s&maxwidth=%d", resURL, opts.MaxWidth)
	}

	if opts.MaxHeight > 0 {
		resURL = fmt.Sprintf("%s&maxheight=%d", resURL, opts.MaxHeight)
	}

	if len(opts.ExtraOpts) > 0 {
		resURL = fmt.Sprintf("%s%s", resURL, buildQueryString(opts.ExtraOpts))
	}

	req, err := http.NewRequest("GET", resURL, nil)

	if err != nil {
		return nil, err
	}

	if len(opts.AcceptLanguage) > 0 {
		req.Header.Add("Accept-Language", opts.AcceptLanguage)
	}

	var resp *http.Response
	if opts.Client != nil {
		resp, err = opts.Client.Do(req)
	} else {
		client := &http.Client{}
		resp, err = client.Do(req)
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return item.parseOembed(opts.URL, resp)
}

// MatchURL tests if given url applies to the endpoint
func (item *Item) MatchURL(url string) bool {
	return item.regex.MatchString(strings.Trim(url, "\r\n"))
}

// NewOembed creates Oembed instance
func NewOembed() *Oembed {
	return &Oembed{}
}

// ParseProviders build oembed endpoint list based on provided json stream
func (o *Oembed) ParseProviders(buf io.Reader) error {
	var providers []Provider

	data, err := ioutil.ReadAll(buf)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &providers)

	if err != nil {
		return err
	}

	var items []*Item

	for _, provider := range providers {
		for _, endpoint := range provider.Endpoints {
			if len(endpoint.Schemes) == 0 {
				endpoint.Schemes = append(endpoint.Schemes, strings.TrimRight(provider.URL, "/")+"/*")
			}
			for _, schema := range endpoint.Schemes {
				or := &Item{ProviderName: provider.Name, ProviderURL: provider.URL}
				or.EndpointURL = o.prepareEndpointURL(endpoint.URL)
				or.regex = o.convertSchemaURL2Regexp(schema)
				items = append(items, or)
			}
		}
	}

	o.items = items

	return nil
}

// FindItem returns Oembed item based on provided url
func (o *Oembed) FindItem(url string) *Item {
	for _, or := range o.items {
		if or.MatchURL(url) {
			return or
		}
	}

	return nil
}

// TODO: add more intelligent parameters parsing
func (o *Oembed) prepareEndpointURL(url string) string {
	url = strings.Replace(url, "{format}", "json", -1)
	url = strings.Replace(url, "/*", "", -1) // hack for Ora TV.. wtf they put in?
	if strings.IndexRune(url, '?') == -1 {
		url += "?format=json&url="
	} else {
		url += "&format=json&url="
	}

	return url
}

func (o *Oembed) convertSchemaURL2Regexp(url string) *regexp.Regexp {
	// domain replacements
	url = strings.Replace(url, "?", "\\?", -1)
	url = su2re1.ReplaceAllString(url, "${1}[^/]%?${2}")
	url = su2re2.ReplaceAllString(url, "${1}.%?${2}")
	url = su2re3.ReplaceAllString(url, "${1}.%")
	url = su2re4.ReplaceAllString(url, "https?://")
	url = strings.Replace(url, "%", "*", -1)
	////
	res, err := regexp.Compile("^" + url + "$")

	if err != nil {
		panic(err)
	}

	return res
}
