package tests

import (
	"ascii-art-web-dockerize/src"
	"testing"
)

func TestIsValidAscii(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"multiple newlines", "hello\n\n\nworld", true},
		{"printable ascii", " !\"#$%&'(){}[]*+,-./09:;<=>?@AZ^_`az|~", true},
		{"low invalid char (below 32)", "hello\x02world", false},
		{"high invalid char (above 126)", "hello\x7fworld", false},
		{"unicode char", "ώ", false},
		{"tab character", "hello\tworld", false},
		{"carriage return", "hello\rworld", false},
		{"backspace", "hello\bworld", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := src.IsValidAscii(tt.input)
			if result != tt.want {
				t.Errorf("IsValidAscii(%q) = %v, want %v", tt.input, result, tt.want)
			}
		})
	}
}
