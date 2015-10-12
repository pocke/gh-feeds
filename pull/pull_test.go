package pull

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPull(t *testing.T) {
	r, err := ioutil.ReadFile("testdata/feeds.xml")
	if err != nil {
		t.Fatal(err)
	}
	f := MockHTTPReq("pocke.private.atom", r)
	defer f()

	uri := "https://github.com/pocke.private.atom?token=tokentokentokentoken"

	_, err = Pull(uri, 1)
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
	r, err := ioutil.ReadFile("testdata/feeds.xml")
	if err != nil {
		t.Fatal(err)
	}
	f := MockHTTPReq("pocke.private.atom", r)
	defer f()

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

type TransportMock struct {
	f func(*http.Request) (*http.Response, error)
}

func (t *TransportMock) RoundTrip(req *http.Request) (*http.Response, error) { return t.f(req) }

var _ http.RoundTripper = &TransportMock{}

type RWCMock struct {
	*bytes.Buffer
}

func (_ *RWCMock) Close() error { return nil }

func MockHTTPReq(uri string, ret []byte) func() {
	bak := http.DefaultTransport
	rt := &TransportMock{
		f: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{}
			body := &RWCMock{
				Buffer: bytes.NewBuffer(ret),
			}
			resp.Body = body
			return resp, nil
		},
	}
	http.DefaultTransport = rt

	return func() { http.DefaultTransport = bak }
}
