package tests

import (
	"ascii-art-web-dockerize/src"
	"strings"
	"testing"
)

func TestParseCLIArgs(t *testing.T) {
	t.Run("single input uses defaults", func(t *testing.T) {
		cfg, err := src.ParseCLIArgs([]string{"hello"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Input != "hello" {
			t.Fatalf("input = %q, want %q", cfg.Input, "hello")
		}
		if cfg.Font != src.Standard {
			t.Fatalf("font = %q, want %q", cfg.Font, src.Standard)
		}
		if cfg.Options.Align != "left" {
			t.Fatalf("align = %q, want %q", cfg.Options.Align, "left")
		}
	})

	t.Run("color substring font", func(t *testing.T) {
		cfg, err := src.ParseCLIArgs([]string{"--color=rgb(255, 0, 0)", "e", "hello", "shadow"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Input != "hello" || cfg.Font != src.Shadow || cfg.Options.SubStr != "e" {
			t.Fatalf("unexpected cfg: %+v", cfg)
		}
		if cfg.Options.ColorHex == nil {
			t.Fatal("expected parsed color, got nil")
		}
		if got := *cfg.Options.ColorHex; got != 0xFF0000 {
			t.Fatalf("color = %#x, want %#x", got, 0xFF0000)
		}
	})

	t.Run("color with positional substring and input uses default font", func(t *testing.T) {
		cfg, err := src.ParseCLIArgs([]string{"--color=red", "e", "hello"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Input != "hello" {
			t.Fatalf("input = %q, want %q", cfg.Input, "hello")
		}
		if cfg.Options.SubStr != "e" {
			t.Fatalf("substr = %q, want %q", cfg.Options.SubStr, "e")
		}
		if cfg.Font != src.Standard {
			t.Fatalf("font = %q, want %q", cfg.Font, src.Standard)
		}
		if cfg.Options.ColorHex == nil || *cfg.Options.ColorHex != 0xFF0000 {
			t.Fatalf("unexpected color value: %+v", cfg.Options.ColorHex)
		}
	})

	t.Run("substring requires color", func(t *testing.T) {
		_, err := src.ParseCLIArgs([]string{"--s", "e", "hello"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		// if !strings.Contains(err.Error(), "substring flag requires --color") {
		// 	t.Fatalf("unexpected error: %v", err)
		// }
	})

	t.Run("invalid font", func(t *testing.T) {
		_, err := src.ParseCLIArgs([]string{"hello", "shadw"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "wrong banner font") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid arg count", func(t *testing.T) {
		_, err := src.ParseCLIArgs([]string{"a", "b", "c", "d"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err != src.ErrUsage {
			t.Fatalf("got %v, want ErrUsage", err)
		}
	})
}
