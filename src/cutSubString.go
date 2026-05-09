package src

import "strings"

type StringColor struct {
	Str  string
	Flag bool
}

// If subStr is empty, the entire string is returned as a single StringColor with Flag: false.
// Returns: a slice of StringColor structs containing the split parts with their color flags
// Example: abcXXXdefg -> [{Str: abc, Flag: false}, {Str: XXX, Flag: true}, {Str: defg, Flag: false}]
func CutSubString(str, subStr string) []StringColor {
	if len(subStr) == 0 {
		return []StringColor{{Str: str, Flag: false}}
	}
	var result []StringColor
	for {
		idx := strings.Index(str, subStr)
		if idx == -1 {
			break
		}
		if idx > 0 {
			result = append(result, StringColor{Str: str[:idx], Flag: false})
		}
		// Add the matched substring with color
		result = append(result, StringColor{Str: subStr, Flag: true})

		str = str[idx+len(subStr):]
	}
	// Remaining text
	if len(str) > 0 {
		result = append(result, StringColor{Str: str, Flag: false})
	}

	return result
}
