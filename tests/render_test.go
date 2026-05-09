package tests

import (
	"ascii-art-web-dockerize/src"
	"strings"
	"testing"
)

func stubFontMap() map[rune][]string {
	makeGlyph := func(ch string) []string {
		return []string{ch, ch, ch, ch, ch, ch, ch, ch}
	}
	return map[rune][]string{
		'A': makeGlyph("A"),
		'B': makeGlyph("B"),
		' ': makeGlyph(" "),
	}
}

func TestRenderInputString(t *testing.T) {
	fontMap := stubFontMap()
	opts := src.BannerOptions{Align: "left"}

	t.Run("escaped newline creates two 8-line blocks", func(t *testing.T) {
		input := strings.ReplaceAll("A\\nB", "\\n", "\n")
		got := src.RenderInputString(input, fontMap, opts, 80)
		lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
		if len(lines) != 16 {
			t.Fatalf("line count = %d, want 16", len(lines))
		}
	})

	t.Run("trailing escaped newline does not create extra block", func(t *testing.T) {
		input := strings.ReplaceAll("A\\n", "\\n", "\n")
		got := src.RenderInputString(input, fontMap, opts, 80)
		lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
		if len(lines) != 8 {
			t.Fatalf("line count = %d, want 8", len(lines))
		}
	})

	t.Run("justify with one word falls back to left layout", func(t *testing.T) {
		left := src.RenderInputString("A", fontMap, src.BannerOptions{Align: "left"}, 2)
		justify := src.RenderInputString("A", fontMap, src.BannerOptions{Align: "justify"}, 2)
		if justify != left {
			t.Fatalf("justify output mismatch\nleft:\n%q\njustify:\n%q", left, justify)
		}
	})

	t.Run("empty input returns empty string", func(t *testing.T) {
		got := src.RenderInputString("", fontMap, opts, 80)
		if got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})

	t.Run("newline only returns newline", func(t *testing.T) {
		got := src.RenderInputString("\n", fontMap, opts, 80)
		if got != "\n" {
			t.Errorf("got %q, want newline", got)
		}
	})

	t.Run("double newline creates eight blank lines between blocks", func(t *testing.T) {
		input := strings.ReplaceAll("A\\n\\nB", "\\n", "\n")
		got := src.RenderInputString(input, fontMap, opts, 80)
		lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
		// A (8 lines) + empty (8 lines) + B (8 lines) = 24 lines
		if len(lines) != 24 {
			t.Fatalf("line count = %d, want 24", len(lines))
		}
	})
}
