package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRender_Page(t *testing.T) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	testRenderer.Renderer = ""
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err == nil {
		t.Fatal(err)
	}
}

func TestRender_GoPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		t.Fatal(err)
	}

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = testRenderer.Page(w, r, "should-not-exist", nil, nil)
	if err == nil {
		t.Fatal(err)
	}
}

func TestRender_NoRender(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		t.Fatal(err)
	}

	testRenderer.Renderer = ""
	testRenderer.RootPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err == nil {
		t.Fatal(err)
	}
}

func TestRender_JetPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		t.Fatal(err)
	}

	testRenderer.Renderer = "jet"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = testRenderer.Page(w, r, "should-not-exist", nil, nil)
	if err == nil {
		t.Fatal(err)
	}
}
