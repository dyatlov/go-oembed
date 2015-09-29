package oembed

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
)

// Info returns information for embedding website
type Info struct {
	Status          int         `json:"-"`
	Type            string      `json:"type"`
	URL             string      `json:"url"`
	ProviderURL     string      `json:"provider_url"`
	ProviderName    string      `json:"provider_name"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Width           json.Number `json:"width"`
	Height          json.Number `json:"height"`
	ThumbnailURL    string      `json:"thumbnail_url"`
	ThumbnailWidth  json.Number `json:"thumbnail_width"`
	ThumbnailHeight json.Number `json:"thumbnail_height"`
	AuthorName      string      `json:"author_name"`
	AuthorURL       string      `json:"author_url"`
	HTML            string      `json:"html"`
}

// NewInfo creater new instance of oembed.Info
func NewInfo() *Info {
	return &Info{}
}

// FillFromJSON fills the structure from provided Oembed JSON
func (info *Info) FillFromJSON(r io.Reader) error {
	data, err := ioutil.ReadAll(r)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &info)

	if err != nil {
		return err
	}

	var t int64
	t, _ = strconv.ParseInt(info.Width.String(), 10, 64)
	info.Width = json.Number(strconv.FormatInt(t, 10))

	t, _ = strconv.ParseInt(info.Height.String(), 10, 64)
	info.Height = json.Number(strconv.FormatInt(t, 10))

	t, _ = strconv.ParseInt(info.ThumbnailWidth.String(), 10, 64)
	info.ThumbnailWidth = json.Number(strconv.FormatInt(t, 10))

	t, _ = strconv.ParseInt(info.ThumbnailHeight.String(), 10, 64)
	info.ThumbnailHeight = json.Number(strconv.FormatInt(t, 10))

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
