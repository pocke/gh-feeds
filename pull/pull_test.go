package pull

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jarcoal/httpmock"
	"github.com/pocke/gh-feeds/db"
)

func TestPull(t *testing.T) {
	MockHTTP()
	defer httpmock.DeactivateAndReset()

	uri := "https://github.com/pocke.private?token=tokentokentokentoken"
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
	MockHTTP()
	defer httpmock.DeactivateAndReset()

	uri := "https://github.com/pocke.private?token=tokentokentokentoken"
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

func TestFeedURI(t *testing.T) {
	MockHTTP()
	defer httpmock.DeactivateAndReset()

	UseTestDB()
	u, err := db.CreateUser(&db.UserParams{
		ID:   1,
		Name: "pocke",
		Auth: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	uri, err := feedURI(u.ID)
	if err != nil {
		t.Fatal(err)
	}

	e := "https://github.com/pocke.private?token=abc123"
	if uri != e {
		t.Errorf("Expected: %s, but got %s", e, uri)
	}
}

func TestDo(t *testing.T) {
	MockHTTP()
	defer httpmock.DeactivateAndReset()

	UseTestDB()
	u, err := db.CreateUser(&db.UserParams{
		ID:   1,
		Name: "pocke",
		Auth: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = Do(u.ID, 1)
	if err != nil {
		t.Error(err)
	}
}

func UseTestDB() {
	s, err := ioutil.ReadFile("../mysql/setup.sql")
	if err != nil {
		panic(err)
	}
	s = []byte(strings.Replace(string(s), "ghfeeds", "ghfeeds_test", -1))
	s = append([]byte("drop database ghfeeds_test;\n"), s...)
	buf := bytes.NewBuffer(s)
	c := exec.Command("mysql", "-uroot")
	c.Stdin = buf
	err = c.Run()
	if err != nil {
		panic(err)
	}

	d, err := sql.Open("mysql", "root:@/ghfeeds_test")
	if err != nil {
		panic(err)
	}
	db.Use(d)
}

func MockHTTP() {
	httpmock.Activate()
	f, err := ioutil.ReadFile("testdata/feeds.xml")
	if err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("GET", "https://github.com/pocke.private",
		httpmock.NewBytesResponder(http.StatusOK, f))

	f, err = ioutil.ReadFile("testdata/feeds_api.json")
	if err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("GET", "https://api.github.com/feeds",
		httpmock.NewBytesResponder(http.StatusOK, f))
}
