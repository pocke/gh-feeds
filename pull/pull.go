package pull

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/tools/blog/atom"
)

type feed struct {
	Entry []entry `xml:"entry"`
}

type entry struct {
	Thumbnail thumbnail `xml:"thumbnail"`
	atom.Entry
}

type thumbnail struct {
	URL string `xml:"url,attr"`
}

func Pull(uri string, page int) (*feed, error) {
	if page < 1 || 10 < page {
		return nil, fmt.Errorf("page should be 1..10, but got %d", page)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &feed{}
	err = xml.NewDecoder(resp.Body).Decode(res)
	return res, err
}
