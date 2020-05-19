package listpage

import (
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"kaldager.com/bysykkel/lib/oslobysykkel"
)

func TestTemplateError(t *testing.T) {
	if err := listpageTmpl.Execute(ioutil.Discard, &pageParams{
		Error: errors.New("oops"),
	}); err != nil {
		t.Fatal(err)
	}
}

func TestTemplateStatic(t *testing.T) {
	src, err := oslobysykkel.NewStaticDataSource("../../testdata/station_status.json", "../../testdata/station_information.json")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	h := Handler{src}
	if err := h.serveHTTP(resp, req); err != nil {
		t.Error(err)
	}
}
