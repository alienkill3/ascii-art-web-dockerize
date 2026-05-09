package tests

import (
	"ascii-art-web-dockerize/src"
	"io"
	"strings"
	"testing"
)

func TestBuildFontMap(t *testing.T) {
	t.Run("builds map from embedded standard font", func(t *testing.T) {
		fontMap, err := src.BuildFontMap(src.GetReader(src.Standard))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(fontMap) != 95 {
			t.Fatalf("len(fontMap) = %d, want 95", len(fontMap))
		}
		if len(fontMap['A']) != 8 {
			t.Fatalf("len(fontMap['A']) = %d, want 8", len(fontMap['A']))
		}
	})

	t.Run("supports CRLF input", func(t *testing.T) {
		raw, err := io.ReadAll(src.GetReader(src.Standard))
		if err != nil {
			t.Fatalf("failed to read standard font: %v", err)
		}
		crlfData := strings.ReplaceAll(string(raw), "\n", "\r\n")
		fontMap, err := src.BuildFontMap(strings.NewReader(crlfData))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(fontMap) != 95 {
			t.Fatalf("len(fontMap) = %d, want 95", len(fontMap))
		}
	})

	t.Run("returns error for truncated font data", func(t *testing.T) {
		truncated := strings.Join([]string{
			"header",
			"a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8",
			"",
		}, "\n")

		_, err := src.BuildFontMap(strings.NewReader(truncated))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "ended prematurely") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
