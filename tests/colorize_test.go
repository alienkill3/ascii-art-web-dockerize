package tests

import (
	"ascii-art-web-dockerize/src"
	"strings"
	"testing"
)

func TestColorizeAndParsers(t *testing.T) {
	t.Run("colorize empty string", func(t *testing.T) {
		if got := src.Colorize("", 0xFF0000); got != "" {
			t.Fatalf("got %q, want empty string", got)
		}
	})

	t.Run("parse hex color", func(t *testing.T) {
		got, err := src.ParseHexColor("#00ff00")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 0x00FF00 {
			t.Fatalf("got %#x, want %#x", got, 0x00FF00)
		}
	})

	t.Run("parse invalid hex color", func(t *testing.T) {
		_, err := src.ParseHexColor("#00gg00")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("parse rgb color", func(t *testing.T) {
		got, err := src.ParseRGBColor("rgb(255, 0, 128)")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 0xFF0080 {
			t.Fatalf("got %#x, want %#x", got, 0xFF0080)
		}
	})

	t.Run("parse invalid rgb color", func(t *testing.T) {
		_, err := src.ParseRGBColor("rgb(300, 0, 0)")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("parse named color is case-insensitive", func(t *testing.T) {
		got, err := src.ParseNamedColor("ReD")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 0xFF0000 {
			t.Fatalf("got %#x, want %#x", got, 0xFF0000)
		}
	})

	t.Run("unsupported named color", func(t *testing.T) {
		_, err := src.ParseNamedColor("rad")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "unsupported color name") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
