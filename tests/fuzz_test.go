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

/*
Go's fuzzer

Seed Corpus:

	It starts with the inputs you provided via f.Add()
	(valid text, edge cases, specific banners).

Execution & Coverage:

	It runs your test function with these seeds and records
	exactly which lines of code are executed (code coverage).

Mutation:

	It then takes a seed and modifies it randomly
	(flipping bits, truncating strings, inserting bytes, etc.).

Feedback Loop:

	It runs the mutated input.
	If this new input causes your code to execute a new path
	(lines of code not reached before), the fuzzer marks this
	input as "interesting" and adds it to the generated corpus.

Crash Detection:

	If any input causes a panic or fails a test assertion
	(like t.Errorf), the fuzzer halts and reports that specific input.

ex. `go test ./tests/ -fuzz=FuzzApiAsciiArtHandler -fuzztime=10s`
*/
func FuzzApiAsciiArtHandler(f *testing.F) {
	// Setup temporary environment with templates directory for the fuzz test
	tmpDir := f.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		f.Fatal(err)
	}
	defer os.Chdir(cwd)

	templatesDir := filepath.Join(tmpDir, "templates")
	if err := os.Mkdir(templatesDir, 0755); err != nil {
		f.Fatal(err)
	}

	// Create a minimal index.html to satisfy the handler's file read
	dummyTemplate := "{{.ErrorMessage}}"
	if err := os.WriteFile(filepath.Join(templatesDir, "index.html"), []byte(dummyTemplate), 0644); err != nil {
		f.Fatal(err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		f.Fatal(err)
	}

	// Add valid cases to the seed corpus
	f.Add("Hello World", "standard")
	f.Add("123456", "shadow")
	f.Add("Testing", "thinkertoy")
	f.Add("", "standard")
	f.Add("Hello\n\nThere", "standard")
	f.Add(" !\"#$%&'(){}[]*+,-./09:;<=>?@AZ^_`az|~", "standard")
	f.Add("hello", "invalid-banner")
	f.Add("hellω", "standard")
	f.Add("hello\x00world", "standard")
	f.Add("Line1\x03Line2", "standard")

	f.Fuzz(func(t *testing.T, text string, banner string) {
		payload, err := json.Marshal(map[string]string{
			"text":   text,
			"banner": banner,
		})
		if err != nil {
			t.Fatalf("failed to marshal request: %v", err)
		}

		req := httptest.NewRequest("POST", "/api/ascii-art", strings.NewReader(string(payload)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		src.ApiAsciiArtHandler(w, req)

		// Replicate handler's validation logic to determine expected status
		// 1. Handler normalizes CRLF to LF
		normalizedText := strings.ReplaceAll(text, "\r\n", "\n")

		// 2. Check conditions that should trigger 400 Bad Request
		// maxInputLength is 2000 in handlers.go
		if len(normalizedText) > 2000 ||
			!src.IsValidAscii(normalizedText) ||
			(banner != "" && !src.IsValidFont(banner)) {
			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 Bad Request for invalid input, got %d. Input: %s, Banner: %q", w.Code, normalizedText, banner)
			}
		} else {
			// 3. Valid inputs should return 200 OK
			if w.Code != http.StatusOK {
				t.Errorf("Expected 200 OK for valid input, got %d. Input: %s, Banner: %q, Body: %s", w.Code, normalizedText, banner, w.Body.String())
			}
		}
	})
}
