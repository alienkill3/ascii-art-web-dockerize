package tests

import (
	"ascii-art-web-dockerize/src"
	"strings"
	"testing"
)

// FuzzParseHexColor tests the hex color parser with random strings.
// It ensures the parser returns an error (and doesn't panic) for invalid inputs.
func FuzzParseHexColor(f *testing.F) {
	// Seed valid and edge cases
	f.Add("#FF0000")
	f.Add("#ffffff")
	f.Add(" #123456 ") // Leading/trailing spaces
	f.Add("invalid")
	f.Add("#123")    // Too short
	f.Add("")        // Empty
	f.Add("#00ff00") // From colorize_test.go
	f.Add("#00gg00") // From colorize_test.go (invalid hex)

	f.Fuzz(func(t *testing.T, color string) {
		_, _ = src.ParseHexColor(color)
	})
}

// FuzzParseRGBColor tests the RGB color parser logic.
// This is critical because it involves splitting strings and parsing integers.
func FuzzParseRGBColor(f *testing.F) {
	f.Add("rgb(255, 0, 0)")
	f.Add("rgb(0,0,0)")
	f.Add("rgb(255,255,255)")
	f.Add(" rgb( 10, 20, 30 ) ")
	f.Add("rgb(256, 0, 0)")
	f.Add("rgb(-1, 0, 0)")
	f.Add("rgb(,,)")
	f.Add("invalid")
	f.Add("rgb(255, 0, 128)")
	f.Add("rgb(300, 0, 0)")

	f.Fuzz(func(t *testing.T, color string) {
		_, _ = src.ParseRGBColor(color)
	})
}

// FuzzParseCLIArgs tests the command line argument parser.
// We simulate CLI args by splitting a random string into a slice.
func FuzzParseCLIArgs(f *testing.F) {
	f.Add("hello standard")
	f.Add("--color=red hello shadow")
	f.Add("--output=file.txt testing")
	f.Add("--align=center --color=#ffffff multi arg test")
	f.Add("hello")
	f.Add("--color=rgb(255,0,0) e hello shadow")
	f.Add("--color=red e hello")
	f.Add("--s e hello")
	f.Add("hello shadw")
	f.Add("a b c d")

	f.Fuzz(func(t *testing.T, input string) {
		// Simulate shell splitting simply by space
		args := strings.Fields(input)
		// We expect errors often, but we must never panic
		_, _ = src.ParseCLIArgs(args)
	})
}
