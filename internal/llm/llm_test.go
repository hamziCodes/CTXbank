package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOllamaClientWithMock(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/tags" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"models":[{"name":"llama3.2"}]}`))
			return
		}

		if r.URL.Path == "/api/generate" {
			w.WriteHeader(http.StatusOK)
			resp := map[string]string{
				"response": "systemPatterns",
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := NewOllamaClient(mockServer.URL, "llama3.2")
	if !client.IsAvailable() {
		t.Errorf("mock server should be reported as available")
	}

	doc, err := client.ClassifyChunk("Architecture and component layout")
	if err != nil {
		t.Fatalf("ClassifyChunk failed: %v", err)
	}

	if doc != DocSystemPatterns {
		t.Errorf("expected DocSystemPatterns, got %s", doc)
	}
}

func TestOllamaClientOfflineFallback(t *testing.T) {
	// Point to unreachable port
	client := NewOllamaClient("http://127.0.0.1:59999", "llama3.2")
	if client.IsAvailable() {
		t.Errorf("offline server should not be available")
	}

	doc, err := client.ClassifyChunk("Some arbitrary note")
	if err == nil {
		t.Errorf("expected ErrLLMUnavailable when offline")
	}
	// Fallback to default productContext
	if doc != DocProductContext {
		t.Errorf("expected fallback to DocProductContext, got %s", doc)
	}
}
