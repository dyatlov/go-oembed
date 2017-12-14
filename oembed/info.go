package oembed

import (
	"encoding/json"
	"strconv"
	"io"

	"github.com/Jeffail/gabs"
)

// Info returns information for embedding website
type Info struct {
	Status          int    `json:"-"`
	Type            string `json:"type,omitempty"`
	CacheAge        uint64 `json:"cache_age,omitempty"`
	URL             string `json:"url,omitempty"`
	ProviderURL     string `json:"provider_url,omitempty"`
	ProviderName    string `json:"provider_name,omitempty"`
	Title           string `json:"title,omitempty"`
	Description     string `json:"description,omitempty"`
	Width           uint64 `json:"width,omitempty"`
	Height          uint64 `json:"height,omitempty"`
	ThumbnailURL    string `json:"thumbnail_url,omitempty"`
	ThumbnailWidth  uint64 `json:"thumbnail_width,omitempty"`
	ThumbnailHeight uint64 `json:"thumbnail_height,omitempty"`
	AuthorName      string `json:"author_name,omitempty"`
	AuthorURL       string `json:"author_url,omitempty"`
	HTML            string `json:"html,omitempty"`
}

// NewInfo creater new instance of oembed.Info
func NewInfo() *Info {
	return &Info{}
}

func jsonUint64(n json.Number) (uint64, error) {
    return strconv.ParseUint(string(n), 10, 64)
}

// FillFromJSON fills the structure from provided Oembed JSON
func (info *Info) FillFromJSON(r io.Reader) error {
	data := json.NewDecoder(r)
	data.UseNumber()

	// We are not using standard json parsing into struct mechanism because it sucks in real life..
	// when you expect a string some sites will return string, some will return null but some will return false
	// when you expect an integer, some will return integer, some will return string but some will return null or false..
	jsonParsed, err := gabs.ParseJSONDecoder(data)

	if err != nil {
		return err
	}

	var strVal string
	var jsonNumberVal json.Number
	var jsonNumber uint64
	var ok bool

	if strVal, ok = jsonParsed.Path("type").Data().(string); ok {
		info.Type = strVal
	}

	if jsonNumberVal, ok = jsonParsed.Path("cache_age").Data().(json.Number); ok {
		if jsonNumber, err = jsonUint64(jsonNumberVal); err != nil {
			return err
		}
		info.CacheAge = jsonNumber
	}

	if strVal, ok = jsonParsed.Path("url").Data().(string); ok {
		info.URL = strVal
	}

	if strVal, ok = jsonParsed.Path("provider_url").Data().(string); ok {
		info.ProviderURL = strVal
	}

	if strVal, ok = jsonParsed.Path("provider_name").Data().(string); ok {
		info.ProviderName = strVal
	}

	if strVal, ok = jsonParsed.Path("title").Data().(string); ok {
		info.Title = strVal
	}

	if strVal, ok = jsonParsed.Path("description").Data().(string); ok {
		info.Description = strVal
	}

	if strVal, ok = jsonParsed.Path("summary").Data().(string); ok {
		info.Description = strVal // stupid nyt oembed uses summary not description
	}

	if strVal, ok = jsonParsed.Path("thumbnail_url").Data().(string); ok {
		info.ThumbnailURL = strVal
	}

	if strVal, ok = jsonParsed.Path("author_name").Data().(string); ok {
		info.AuthorName = strVal
	}

	if strVal, ok = jsonParsed.Path("author_url").Data().(string); ok {
		info.AuthorURL = strVal
	}

	if strVal, ok = jsonParsed.Path("html").Data().(string); ok {
		info.HTML = strVal
	}

	if jsonNumberVal, ok = jsonParsed.Path("width").Data().(json.Number); ok {
		if jsonNumber, err = jsonUint64(jsonNumberVal); err != nil {
			return err
		}
		info.Width = jsonNumber
	}

	if jsonNumberVal, ok = jsonParsed.Path("height").Data().(json.Number); ok {
		if jsonNumber, err = jsonUint64(jsonNumberVal); err != nil {
			return err
		}
		info.Height = jsonNumber
	}

	if jsonNumberVal, ok = jsonParsed.Path("thumbnail_width").Data().(json.Number); ok {
		if jsonNumber, err = jsonUint64(jsonNumberVal); err != nil {
			return err
		}
		info.ThumbnailWidth = jsonNumber
	}

	if jsonNumberVal, ok = jsonParsed.Path("thumbnail_height").Data().(json.Number); ok {
		if jsonNumber, err = jsonUint64(jsonNumberVal); err != nil {
			return err
		}
		info.ThumbnailHeight = jsonNumber
	}

	return nil
}

// MergeWith adds missing data to the structure from another oembed.Info structure
func (info *Info) MergeWith(base *Info) {
	if len(info.Type) == 0 {
		info.Type = base.Type
	}
	if len(info.URL) == 0 {
		info.URL = base.URL
	}
	if len(info.ProviderURL) == 0 {
		info.ProviderURL = base.ProviderURL
	}
	if len(info.ProviderName) == 0 {
		info.ProviderName = base.ProviderName
	}
	if len(info.Title) == 0 {
		info.Title = base.Title
	}
	if len(info.Description) == 0 {
		info.Description = base.Description
	}
	if len(info.ThumbnailURL) == 0 {
		info.ThumbnailURL = base.ThumbnailURL
		info.ThumbnailWidth = base.ThumbnailWidth
		info.ThumbnailHeight = base.ThumbnailHeight
	}
}

// ToJSON a simple wrapper around json.Marshal
func (info *Info) ToJSON() ([]byte, error) {
	return json.Marshal(info)
}

// String return json representation of structure, or error string
func (info *Info) String() string {
	data, err := info.ToJSON()

	if err != nil {
		return err.Error()
	}

	return string(data[:])
}
