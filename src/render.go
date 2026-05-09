package src

import (
	"regexp"
	"strings"
)

// ansiRegexp matches ANSI escape sequences used for colorization.
var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripAnsi removes ANSI escape codes from the input string.
func stripAnsi(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// visibleLen returns the character length of the string after stripping
// ANSI escape codes, which are not counted towards terminal width.
func visibleLen(s string) int {
	return len(stripAnsi(s))
}

type BannerOptions struct {
	ColorHex *int
	SubStr   string
	Align    string
}

func bannerBlockWithOptions(inputString string, fontMap map[rune][]string, options BannerOptions, width uint16) string {
	parts := []StringColor{{Str: inputString, Flag: true}}
	if options.ColorHex == nil {
		parts[0].Flag = false
	} else if options.SubStr != "" {
		parts = CutSubString(inputString, options.SubStr)
	}

	var builder strings.Builder
	for line := 0; line < 8; line++ {
		for _, part := range parts {
			for _, rn := range part.Str {
				ascii, exists := fontMap[rn]
				if !exists {
					continue
				}

				glyph := ascii[line]

				if part.Flag && options.ColorHex != nil {
					builder.WriteString(Colorize(glyph, *options.ColorHex))
					continue
				}
				builder.WriteString(glyph)

			}
		}
		builder.WriteByte('\n')
	}
	var newBuilder strings.Builder
	switch options.Align {
	case "right":
		lines := strings.Split(builder.String(), "\n")
		for i, ln := range lines {
			// skip the trailing empty element caused by the final newline
			if i == len(lines)-1 && ln == "" {
				continue
			}
			vis := visibleLen(ln)
			pad := int(width) - vis
			if pad < 0 {
				pad = 0
			}
			newBuilder.WriteString(strings.Repeat(" ", pad) + ln + "\n")
		}
	case "center":
		lines := strings.Split(builder.String(), "\n")
		for i, ln := range lines {
			if i == len(lines)-1 && ln == "" {
				continue
			}
			vis := visibleLen(ln)
			pad := (int(width) - vis) / 2
			if pad < 0 {
				pad = 0
			}
			newBuilder.WriteString(strings.Repeat(" ", pad) + ln + "\n")
		}
	case "justify":
		cleanInput := strings.TrimSpace(inputString)
		words := strings.Fields(cleanInput)
		if len(words) < 2 {
			newBuilder.WriteString(builder.String())
			break
		}

		contentWidth := 0
		for _, part := range parts {
			for _, rn := range part.Str {
				if rn != ' ' {
					if ascii, ok := fontMap[rn]; ok {
						contentWidth += len(ascii[0])
					}
				}
			}
		}

		totalSpaceWidth := int(width) - contentWidth
		numGaps := len(words) - 1
		if totalSpaceWidth < numGaps {
			newBuilder.WriteString(builder.String())
			break
		}

		spacesPerGap := totalSpaceWidth / numGaps
		extraSpaces := totalSpaceWidth % numGaps

		for line := 0; line < 8; line++ {
			currentGap := 0
			seenFirstWord := false
			lastWasSpace := false
			for _, part := range parts {
				for _, rn := range part.Str {
					if rn != ' ' {
						seenFirstWord = true
						lastWasSpace = false
						if ascii, ok := fontMap[rn]; ok {
							glyph := ascii[line]
							if part.Flag && options.ColorHex != nil {
								newBuilder.WriteString(Colorize(glyph, *options.ColorHex))
							} else {
								newBuilder.WriteString(glyph)
							}
						}
					} else {
						if seenFirstWord && !lastWasSpace && currentGap < numGaps {
							w := spacesPerGap
							if currentGap < extraSpaces {
								w++
							}
							newBuilder.WriteString(strings.Repeat(" ", w))
							currentGap++
						}
						lastWasSpace = true
					}
				}
			}
			newBuilder.WriteByte('\n')
		}
	default: //left
		newBuilder = builder
	}

	// return builder.String()
	return newBuilder.String()
}

// func PrintBannerWithOptions(inputString string, fontMap map[rune][]string, options BannerOptions, width uint16) {
// 	fmt.Print(bannerBlockWithOptions(inputString, fontMap, options, width))
// 	fmt.Println("test")
// }

func RenderInputString(input string, fontMap map[rune][]string, options BannerOptions, width uint16) string {
	// Replace escaped newlines with actual newlines
	input = strings.ReplaceAll(input, "\\n", "\n")
	parts := strings.Split(input, "\n")
	var builder strings.Builder

	// Ignore trailing "\n" (it should not create an extra empty output block)
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}

	for _, p := range parts {
		if p == "" {
			// Empty lines render as 8 blank lines (matching character height)
			// But only if there's non-empty content somewhere else (not just a lone newline)
			if len(parts) > 1 {
				// Generate 8 empty lines
				for j := 0; j < 8; j++ {
					builder.WriteByte('\n')
				}
			} else {
				// Lone newline or empty input - just output a newline
				builder.WriteByte('\n')
			}
			continue
		}
		builder.WriteString(bannerBlockWithOptions(p, fontMap, options, width))
	}
	return builder.String()
}

// RenderInputStringWeb renders input string for web interface without escaping \n sequences
func RenderInputStringWeb(input string, fontMap map[rune][]string, options BannerOptions, width uint16) string {
	// For web interface, don't replace \n sequences - treat them as literal characters
	parts := strings.Split(input, "\n")
	var builder strings.Builder

	// Ignore trailing "\n" (it should not create an extra empty output block)
	if len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}

	for _, p := range parts {
		if p == "" {
			// Empty lines render as 8 blank lines (matching character height)
			// But only if there's non-empty content somewhere else (not just a lone newline)
			if len(parts) > 1 {
				// Generate 8 empty lines
				for j := 0; j < 8; j++ {
					builder.WriteByte('\n')
				}
			} else {
				// Lone newline or empty input - just output a newline
				builder.WriteByte('\n')
			}
			continue
		}
		builder.WriteString(bannerBlockWithOptions(p, fontMap, options, width))
	}
	return builder.String()
}

// func RenderInput(input string, fontMap map[rune][]string, options BannerOptions, width uint16) {
// 	fmt.Printf("%*s", width, RenderInputString(input, fontMap, options, width))
// 	fmt.Println("test")
// }
