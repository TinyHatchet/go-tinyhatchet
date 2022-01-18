package mouseion_test

import (
	"io"
	"net/http"
	"testing"

	tinyhatchet "github.com/tinyhatchet/go-tinyhatchet"
)

type MockClient struct {
	post func(url string, contentType string, body io.Reader) (*http.Response, error)
}

func (m *MockClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return m.post(url, contentType, body)
}

func TestArgsToTags(t *testing.T) {
	logger := tinyhatchet.Logger{
		LogErrors:   false,
		DefaultTags: []string{"default tag"},
		AutoTagger: func(defaultTags []string, arg interface{}) []string {
			return []string{"err", defaultTags[0]}
		},
	}
	result := logger.ArgsToTags("test")
	if result == nil {
		t.Fatal("got a nil tag list")
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(result))
	}
}
