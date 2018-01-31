package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/baskeboler/auth"
	"github.com/go-kit/kit/log"
)

func TestMain(m *testing.M) {

	auth.ResetDb()
	os.Exit(m.Run())
}
func TestGenerate(t *testing.T) {
	s, e := NewService()
	if e != nil {
		t.Error(e)
	}

	logger := log.NewNopLogger()
	router := MakeHTTPRouter(s, logger)
	recorder := httptest.NewRecorder()

	b := strings.NewReader(`{
		"name": "user",
		"password": "user"	
	}`)
	// defer b.Close()
	req := httptest.NewRequest("POST", buildUrl("/api/login"), b)
	// req.Body = ioutil.NopCloser(b)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fail()
	}
	t.Log(recorder.Body)
	var keyResp loginResponse

	json.NewDecoder(recorder.Body).Decode(&keyResp)
	t.Log(keyResp)
	req = httptest.NewRequest("GET", buildUrl("/api/puzzle?width=10&height=10&words=5"), nil)
	req.Header.Set("Authorization", keyResp.Key.Key)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fail()
	}
}
func TestLogin(t *testing.T) {
	s, e := NewService()
	if e != nil {
		t.Error(e)
	}
	logger := log.NewLogfmtLogger(os.Stdout)
	router := MakeHTTPRouter(s, logger)
	recorder := httptest.NewRecorder()

	b := strings.NewReader(`{
		"name": "user",
		"password": "user"	
	}`)
	// defer b.Close()
	req := httptest.NewRequest("POST", buildUrl("/api/login"), b)
	// req.Body = ioutil.NopCloser(b)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fail()
	}
	t.Log(recorder.Body)
}
func buildUrl(path string) string {

	return fmt.Sprintf("http://localhost:8080%s", path)
}
