package tests

import (
	"regexp"
	"strings"
	"testing"

	"ascii-art-web-dockerize/src"
)

// stripAnsi is a copy of the helper used by render.go. It is only used
// in tests to verify that the visible width of a line matches expected
// values when the line contains color escape sequences.
func stripAnsi(s string) string {
	var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansi.ReplaceAllString(s, "")
}

func TestColorAlignment(t *testing.T) {
	fontMap := map[rune][]string{
		'A': {"X", "X", "X", "X", "X", "X", "X", "X"},
	}

	hex := 0xFF0000
	colorized := src.Colorize("X", hex)

	widths := []uint16{5, 6}
	aligns := []string{"right", "center", "left"}

	for _, width := range widths {
		for _, align := range aligns {
			opts := src.BannerOptions{ColorHex: &hex, Align: align}
			got := src.RenderInputString("A", fontMap, opts, width)

			lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
			if len(lines) != 8 {
				t.Fatalf("expected 8 lines, got %d", len(lines))
			}

			for _, ln := range lines {
				vis := len(stripAnsi(ln))
				if align == "right" && vis != int(width) {
					t.Errorf("right align width %d: got visible len %d line %q",
						width, vis, ln)
				}
				if align == "center" && vis > int(width) {
					t.Errorf("center align overflow width %d: got visible len %d line %q",
						width, vis, ln)
				}
			}

			switch align {
			case "right":
				pad := int(width) - 1
				if pad < 0 {
					pad = 0
				}
				expectLine := strings.Repeat(" ", pad) + colorized
				for _, ln := range lines {
					if ln != expectLine {
						t.Errorf("right align mismatch (width %d): expected %q got %q",
							width, expectLine, ln)
					}
				}
			case "center":
				pad := (int(width) - 1) / 2
				if pad < 0 {
					pad = 0
				}
				expectLine := strings.Repeat(" ", pad) + colorized
				for _, ln := range lines {
					if ln != expectLine {
						t.Errorf("center align mismatch (width %d): expected %q got %q",
							width, expectLine, ln)
					}
				}
			case "left":
				expectLine := colorized
				for _, ln := range lines {
					if ln != expectLine {
						t.Errorf("left align mismatch: expected %q got %q",
							expectLine, ln)
					}
				}
			}
		}
	}

	// also sanity check justify aligns to width when coloring is enabled
	{
		opts := src.BannerOptions{ColorHex: &hex, Align: "justify"}
		got := src.RenderInputString("A A", fontMap, opts, 5)
		lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
		for _, ln := range lines {
			vis := len(stripAnsi(ln))
			if vis != 5 {
				t.Errorf("justify align visible width %d, want 5, line %q", vis, ln)
			}
		}

	}
}
