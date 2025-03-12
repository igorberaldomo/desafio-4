package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestFindTempHandler(t *testing.T) {
	s:= httptest.NewServer(http.HandlerFunc(FindTempHandler))

	req,err := http.NewRequest(http.MethodGet, s.URL+"?cep=78050040", nil)
	if err != nil {
		t.Fatal(err)
	}

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

	wantRegexp := regexp.MustCompile(`^{"temp_C":\d+(\.\d+)?,"temp_F":\d+(\.\d+)?,"temp_K":\d+(\.\d+)?}$`)

	sbody :=string(body)
	sbody = strings.Trim(sbody, " \n")

	if !wantRegexp.MatchString(sbody) {
		t.Fatalf("want string is not valid '%s' ", sbody)
	}
}
