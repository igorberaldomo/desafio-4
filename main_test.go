package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindTempHandler(t *testing.T) {
	s:= httptest.NewServer(http.HandlerFunc(FindTempHandler))

	req := httptest.NewRequest(http.MethodGet, s.URL+"/?cep=78050040", nil)

	client := &http.Client{}
	res,err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	want := `{"temp_C":28.5,"temp_F":28.5,"temp_K":28.5}`

	if string(body) != want {
		t.Errorf("want %q, got %q", want, string(body))
	}
}
