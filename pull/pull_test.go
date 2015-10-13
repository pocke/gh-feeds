package pull

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func MockFeeds() {
	httpmock.Activate()
	f, err := ioutil.ReadFile("testdata/feeds.xml")
	if err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("GET", "https://github.com/pocke.private.atom",
		httpmock.NewBytesResponder(http.StatusOK, f))
}

func TestPull(t *testing.T) {
	MockFeeds()
	defer httpmock.DeactivateAndReset()

	uri := "https://github.com/pocke.private.atom?token=tokentokentokentoken"
	_, err := Pull(uri, 1)
	if err != nil {
		t.Fatal(err)
	}

	// When page is out of range
	_, err = Pull(uri, -1)
	if err == nil {
		t.Error("should be error when page is out of range, but got nil")
	}
	_, err = Pull(uri, 0)
	if err == nil {
		t.Error("should be error when page is out of range, but got nil")
	}
	_, err = Pull(uri, 11)
	if err == nil {
		t.Error("should be error when page is out of range, but got nil")
	}
}

func TestTransform(t *testing.T) {
	MockFeeds()
	defer httpmock.DeactivateAndReset()

	uri := "https://github.com/pocke.private.atom?token=tokentokentokentoken"
	resp, err := Pull(uri, 1)
	if err != nil {
		t.Fatal(err)
	}

	evs, err := transform(resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(evs) != len(resp.Entry) {
		t.Fatalf("%d != %d", len(evs), len(resp.Entry))
	}
}
