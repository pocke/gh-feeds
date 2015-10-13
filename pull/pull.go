package pull

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/go-github/github"
	"github.com/joeshaw/iso8601"
	"github.com/pocke/gh-feeds/db"

	"golang.org/x/oauth2"
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

func timeStrToTime(s atom.TimeStr) (time.Time, error) {
	// XXX: これでいい?
	return time.Parse(iso8601.Format+"Z", string(s))
}

func transform(f *feed) ([]db.Event, error) {
	events := make([]db.Event, 0, len(f.Entry))
	for _, e := range f.Entry {
		t, err := timeStrToTime(e.Published)
		if err != nil {
			return nil, err
		}

		var url string
		if len(e.Link) != 0 {
			url = e.Link[0].Href
		}

		ev := db.Event{
			PublishedAt: t,
			Type:        e.ID, // TODO:
			HTML:        e.Content.Body,
			AuthorName:  e.Author.Name,
			UserId:      1, // TODO:
			URL:         url,
			ImageURL:    e.Thumbnail.URL,
		}
		events = append(events, ev)
	}
	return events, nil
}

func feedURI(user_id int) (string, error) {
	u, err := db.User{}.Find(user_id)
	if err != nil {
		return "", err
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: u.Auth},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	c := github.NewClient(tc)
	fe, _, err := c.Feeds.ListFeeds()
	if err != nil {
		return "", err
	}
	return fe.Links.CurrentUser.Href, nil
}

// TODO: URIをもらうんじゃなく、user_idを貰ってAPIを叩いてあれをあれする
func Do(user_id int, page int) error {
	uri, err := feedURI(user_id)
	if err != nil {
		return err
	}

	f, err := Pull(uri, page)
	if err != nil {
		return err
	}
	evs, err := transform(f)
	if err != nil {
		return err
	}
	// TODO: バルクインサート
	for _, e := range evs {
		_, err := e.Save()
		if err != nil {
			return err
		}
	}
	return nil
}
