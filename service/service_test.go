package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
)

func TestGenerate(t *testing.T) {
	s, e := NewService()
	if e != nil {
		t.Error(e)
	}
	logger := log.NewNopLogger()
	router := MakeHTTPRouter(s, logger)
	recorder := httptest.NewRecorder()

	req := httptest.NewRequest("GET", buildUrl("/api/puzzle?width=10&height=10&words=5"), nil)

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fail()
	}
}

func buildUrl(path string) string {

	return fmt.Sprintf("http://localhost:8080%s", path)
}
