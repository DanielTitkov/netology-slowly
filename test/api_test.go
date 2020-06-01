package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DanielTitkov/netology-slowly/internal/app"
	"github.com/DanielTitkov/netology-slowly/internal/configs"
	"github.com/DanielTitkov/netology-slowly/internal/middleware"
	"github.com/go-chi/chi"
)

const slowURL = "/api/slow"
const maxSlowTimeout = 50

func makeServer() *httptest.Server {
	cfg := configs.Config{
		MaxSlowTimeout: maxSlowTimeout,
	}
	handler := app.NewApp(cfg)
	router := chi.NewRouter()
	router.With(middleware.NewTimeout(cfg)).Post(slowURL, handler.SlowHandler)
	return httptest.NewServer(router)
}

func Test404(t *testing.T) {
	expected := http.StatusNotFound
	server := makeServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/some/url")
	if err != nil {
		t.Fatal("failed to get response", err)
	}
	if resp.StatusCode != expected {
		t.Fatalf("expected %d status, got %d", expected, resp.StatusCode)
	}
}

func TestHandleSlowGet(t *testing.T) {
	expected := http.StatusMethodNotAllowed
	server := makeServer()
	defer server.Close()

	resp, err := http.Get(server.URL + slowURL)
	if err != nil {
		t.Fatal("failed to get response", err)
	}
	if resp.StatusCode != expected {
		t.Fatalf("expected %d status, got %d", expected, resp.StatusCode)
	}
}

func TestHandleSlowPost(t *testing.T) {
	tolerance := time.Duration(5) * time.Millisecond
	d := 20 // delta in ms
	type testCase struct {
		JSON        string
		expCode     int
		expBody     string
		expRespTime time.Duration
	}
	testCases := []testCase{
		{fmt.Sprintf(`{"timeout": %d}`, maxSlowTimeout-d), 200, `{"status":"ok"}`, time.Duration(maxSlowTimeout-d) * time.Millisecond},
		{fmt.Sprintf(`{"timeout": %d}`, maxSlowTimeout+d), 400, `{"error":"timeout too long"}`, time.Duration(maxSlowTimeout) * time.Millisecond},
	}

	server := makeServer()
	defer server.Close()
	for i, tc := range testCases {
		start := time.Now()
		resp, err := http.Post(server.URL+slowURL, "application/json;charset=utf-8", bytes.NewBuffer([]byte(tc.JSON)))
		if err != nil {
			t.Fatalf("failed to get response on testcase %d with err %s", i, err)
		}
		respTime := time.Since(start)
		if resp.StatusCode != tc.expCode {
			t.Fatalf("on testcase %d expected %d status, got %d", i, tc.expCode, resp.StatusCode)
		}
		if body, _ := ioutil.ReadAll(resp.Body); strings.Trim(string(body), " \n") != tc.expBody {
			t.Fatalf("on testcase %d expected '%s' response, got '%s'", i, tc.expBody, string(body))
		}
		if respTime > tc.expRespTime+tolerance || respTime < tc.expRespTime-tolerance {
			t.Fatalf("on testcase %d expected request to be completed in %s +/- %s, got %s", i, tc.expRespTime, tolerance, respTime)
		}
	}
}
