package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/handson-go/chap6/complex_server/config"
)

func TestApiHandler(t *testing.T) {
	r := httptest.NewRequest("GET", "/api", nil)
	w := httptest.NewRecorder()

	b := new(bytes.Buffer)
	c := config.InitConfig(b)

	apiHandler(w, r, c)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected response status: %v, Got: %v\n", http.StatusOK, resp.StatusCode)
	}
	expectedResponsBody := "Hello, world!"

	if string(body) != expectedResponsBody {
		t.Errorf("Expected response: %s, Got: %s\n", expectedResponsBody, string(body))
	}
}
