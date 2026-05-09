package src

func IsValidAscii(inputString string) bool {
	for _, rn := range inputString {
		if (rn < 32 || rn > 126) && rn != '\n' {
			// we found non ascii character
			return false
		}
	}
	return true
}
