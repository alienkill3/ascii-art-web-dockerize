package tests

import (
	"ascii-art-web-dockerize/src"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandlers(t *testing.T) {
	// Setup temporary environment with templates directory
	tmpDir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd) // Restore original working directory

	// Create templates directory
	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a dummy index.html
	// We include placeholders for the fields to verify they are populated correctly
	dummyTemplate := `{{.Input}}|{{.Banner}}|{{.Result}}|{{.ErrorMessage}}`
	err = os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(dummyTemplate), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Switch to temp dir so handlers can find "templates/index.html"
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	t.Run("HomeHandler", func(t *testing.T) {
		t.Run("GET / returns 200 and renders template", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			src.HomeHandler(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
			}
			// Check if template rendered default banner and empty error
			if !strings.Contains(w.Body.String(), "|standard||") {
				t.Errorf("unexpected body: %q", w.Body.String())
			}
		})

		t.Run("POST / returns 400", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()

			src.HomeHandler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
			}
			if !strings.Contains(w.Body.String(), "bad request: GET is required") {
				t.Errorf("unexpected body: %q", w.Body.String())
			}
		})

		t.Run("Invalid path returns 404", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
			w := httptest.NewRecorder()

			src.HomeHandler(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
			}
		})

		t.Run("500 Internal Server Error on broken template", func(t *testing.T) {
			// Corrupt the template to force a parsing error
			badTemplate := "{{"
			if err := os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(badTemplate), 0644); err != nil {
				t.Fatal(err)
			}
			defer func() {
				// Restore the valid template
				os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(dummyTemplate), 0644)
			}()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			src.HomeHandler(w, req)

			if w.Code != http.StatusInternalServerError {
				t.Errorf("got status %d, want %d", w.Code, http.StatusInternalServerError)
			}
		})

		t.Run("Template missing returns 404", func(t *testing.T) {
			// Remove the template to force not found error
			templatePath := filepath.Join(templatesDir, "index.html")
			if err := os.Remove(templatePath); err != nil {
				t.Fatal(err)
			}
			defer func() {
				// Restore the valid template
				if err := os.WriteFile(templatePath, []byte(dummyTemplate), 0644); err != nil {
					t.Errorf("failed to restore template: %v", err)
				}
			}()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			src.HomeHandler(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("got status %d, want %d", w.Code, http.StatusNotFound)
			}
		})
	})

	t.Run("ApiAsciiArtHandler", func(t *testing.T) {
		t.Run("GET returns 405", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/ascii-art", nil)
			w := httptest.NewRecorder()

			src.ApiAsciiArtHandler(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("got status %d, want %d", w.Code, http.StatusMethodNotAllowed)
			}
			if !strings.Contains(w.Body.String(), "method must be POST") {
				t.Errorf("unexpected body: %q", w.Body.String())
			}
		})

		t.Run("POST non-ascii input", func(t *testing.T) {
			payload := `{"text":"hellω","banner":"standard"}`
			req := httptest.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			src.ApiAsciiArtHandler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
			}
			if !strings.Contains(w.Body.String(), "input contains non-ASCII characters") {
				t.Errorf("unexpected body: %q", w.Body.String())
			}
		})

		t.Run("POST 400 invalid banner", func(t *testing.T) {
			payload := `{"text":"hello","banner":"invalid-banner"}`
			req := httptest.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			src.ApiAsciiArtHandler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
			}
			if !strings.Contains(w.Body.String(), "invalid banner") {
				t.Errorf("unexpected body: %q", w.Body.String())
			}
		})

		t.Run("API POST json ascii art ok", func(t *testing.T) {
			payload := `{"text":"hi"}`
			req := httptest.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			src.ApiAsciiArtHandler(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("got status %d, want %d", w.Code, http.StatusOK)
			}
			if got := w.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("got Content-Type %q, want %q", got, "application/json")
			}

			var resp struct {
				Result string `json:"result"`
			}
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed decode json: %v", err)
			}
			if resp.Result == "" {
				t.Fatalf("expected non-empty result")
			}
		})

		t.Run("API POST bad json returns 400", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/ascii-art", strings.NewReader("{bad"))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			src.ApiAsciiArtHandler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("got status %d, want %d", w.Code, http.StatusBadRequest)
			}

			var resp struct {
				Error string `json:"error"`
			}
			if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
				t.Fatalf("failed decode json: %v", err)
			}
			if resp.Error != "invalid JSON payload" {
				t.Fatalf("got error %q, want %q", resp.Error, "invalid JSON payload")
			}
		})
	})
}
