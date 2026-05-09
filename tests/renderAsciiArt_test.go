package tests

import (
	"ascii-art-web-dockerize/src"
	"os"
	"path/filepath"
	"testing"
)

func TestRenderASCIIArt(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		banner     string
		goldenFile string
	}{
		{
			name:       "Hello There (Standard)",
			input:      "Hello There",
			banner:     "standard",
			goldenFile: "hello_there.golden",
		},
		{
			name:       "Hello Shadow",
			input:      "hello",
			banner:     "shadow",
			goldenFile: "hello_shadow.golden",
		},
		{
			name:       "Hello Thinkertoy",
			input:      "hello",
			banner:     "thinkertoy",
			goldenFile: "hello_thinkertoy.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := src.RenderASCIIArt(tt.input, tt.banner)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			content, err := os.ReadFile(filepath.Join("..", "testfiles", tt.goldenFile))
			if err != nil {
				t.Fatalf("failed to read golden file: %v", err)
			}

			want := string(content)
			if got != want {
				t.Errorf("result mismatch for %s\ngot:\n%q\nwant:\n%q", tt.name, got, want)
			}
		})
	}
}
