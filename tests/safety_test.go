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

// TestInputSafety verifies that the server handles malicious or extreme inputs gracefully
// without crashing (panicking) or exhausting resources.
func TestInputSafety(t *testing.T) {
	// Setup mock environment with templates to ensure handlers can run
	tmpDir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Minimal template to satisfy handler logic
	dummyTemplate := "{{.ErrorMessage}}"
	if err := os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(dummyTemplate), 0644); err != nil {
		t.Fatal(err)
	}

	// Change working directory so handlers find the "templates" folder
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		input          string
		expectedStatus int
	}{
		{
			name:           "Buffer Overflow Attempt (Oversized)",
			input:          strings.Repeat("A", 3000), // Exceeds 2000 char limit
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Null Byte Injection",
			input:          "Hello\x00World", // Invalid ASCII (often breaks C-strings or parsers)
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Control Characters",
			input:          "Line1\x03Line2", // Non-printable ASCII
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recover from potential panics to report them as test failures
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("CRITICAL: Server panicked on %s: %v", tt.name, r)
				}
			}()

			payload, err := json.Marshal(map[string]string{
				"text":   tt.input,
				"banner": "standard",
			})
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}

			req := httptest.NewRequest("POST", "/api/ascii-art", strings.NewReader(string(payload)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			src.ApiAsciiArtHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d for %s, got %d", tt.expectedStatus, tt.name, w.Code)
			}
		})
	}
}
