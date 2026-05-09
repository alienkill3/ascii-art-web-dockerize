package src

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Colorize applies RGB color formatting to a string using ANSI escape codes.
// for 24-bit true color terminal output.
// Example: colorized := Colorize("Hello", 0xFF0000) // Red color
func Colorize(str string, colorHex int) string {
	if len(str) == 0 {
		return ""
	}
	r := (colorHex >> 16) & 0xFF
	g := (colorHex >> 8) & 0xFF
	b := colorHex & 0xFF
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, str)
}

// ParseHexColor parses a hex RGB string to integer.
// Accepts only "#RRGGBB" format.
func ParseHexColor(color string) (int, error) {
	clean := strings.TrimSpace(color)
	if !strings.HasPrefix(clean, "#") || len(clean) != 7 {
		return 0, errors.New("invalid HexColor format: use #RRGGBB\n")
	}

	val, err := strconv.ParseInt(clean[1:], 16, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid hex color (Valid hex digits are 0-9, A-F)\nExample: #FF0000 (red)\n")
	}

	return int(val), nil
}

// ParseRGBColor parses an RGB color string in format "rgb(R, G, B)".
func ParseRGBColor(color string) (int, error) {
	clean := strings.TrimSpace(color)
	if !strings.HasPrefix(clean, "rgb(") || !strings.HasSuffix(clean, ")") {
		return 0, fmt.Errorf("invalid RGB format: use rgb'(R, G, B)' where R, G, B are 0-255\n")
	}

	// Extract the values inside rgb(...)
	inside := strings.TrimPrefix(clean, "rgb(")
	inside = strings.TrimSuffix(inside, ")")

	parts := strings.Split(inside, ",")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid RGB format: use rgb'(R, G, B)' where R, G, B are 0-255\n")
	}

	var r, g, b int
	for i, part := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || val < 0 || val > 255 {
			return 0, errors.New("invalid RGB value: each component must be 0-255\n")
		}
		switch i {
		case 0:
			r = val
		case 1:
			g = val
		case 2:
			b = val
		}
	}

	return (r << 16) | (g << 8) | b, nil
}

// ParseNamedColor parses a color name to integer RGB value.
// Supported colors: red, blue, green, purple, cyan, yellow
func ParseNamedColor(color string) (int, error) {
	clean := strings.ToLower(strings.TrimSpace(color))

	colorMap := map[string]int{
		"red":    0xFF0000, // (255, 0, 0)
		"blue":   0x0000FF, // (0, 0, 255)
		"green":  0x00FF00, // (0, 255, 0)
		"purple": 0x800080, // (128, 0, 128)
		"cyan":   0x00FFFF, // (0, 255, 255)
		"yellow": 0xFFFF00, // (255, 255, 0)
		"orange": 0xFFA500, // (255, 165, 0)
		"pink":   0xFFC0CB, // (255, 192, 203)

	}

	if val, ok := colorMap[clean]; ok {
		return val, nil
	}

	return 0, fmt.Errorf("unsupported color name: try using a hex code (#RRGGBB) or RGB format (rgb'(R, G, B)')\n")
}
