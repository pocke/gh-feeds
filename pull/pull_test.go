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

	resp, err := Pull("https://github.com/pocke.private.atom?token=tokentokentokentoken", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp.Entry[0].Thumbnail)
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
