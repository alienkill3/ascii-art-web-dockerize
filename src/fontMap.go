package src

import (
	"fmt"
	"io"
	"strings"
)

func BuildFontMap(fontByte io.Reader) (map[rune][]string, error) {
	// Read the entire font into memory at once
	data, err := io.ReadAll(fontByte)
	if err != nil {
		return nil, fmt.Errorf("failed to read font data: %w", err)
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	fontMap := make(map[rune][]string, 95) // printable ASCII
	currentChar := rune(32)

	// skip header line if present
	lineIndex := 1

	for lineIndex+7 < len(lines) {
		charLines := lines[lineIndex : lineIndex+8] // 8 lines per character
		fontMap[currentChar] = charLines
		currentChar++
		lineIndex += 9 // 8 lines + 1 empty line separator
	}

	// sanity check
	if currentChar <= 126 {
		return fontMap, fmt.Errorf("font file ended prematurely; last char = %c", currentChar-1)
	}

	return fontMap, nil
}
