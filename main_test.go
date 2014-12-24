package main

import (
	"github.com/masayukioguni/go-lgtm-front/config"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func NewMockServer() *Front {
	c, _ := config.NewConfig(".env_test")
	f := &Front{
		config: c,
	}
	return f
}

func TestMain_Index(t *testing.T) {
	f := NewMockServer()
	ts := httptest.NewServer(http.HandlerFunc(f.Index))
	defer ts.Close()

	res, err := http.Get(ts.URL)

	if err != nil {
		t.Errorf("TestMain_Index by http.Get() returned %+v", err)
	}

	wantStatusCode := 200
	if !reflect.DeepEqual(res.StatusCode, wantStatusCode) {
		t.Errorf("TestMain_Index Response Code returned %+v, want %+v", res.StatusCode, wantStatusCode)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Errorf("TestMain_Index by ioutil.ReadAll() returned %+v", err)
	}

}
