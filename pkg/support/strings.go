package support

import "bytes"

var specials = []rune{'.', '\'', '(', ')', '-', '+'}

func AddSlashes(str string) string {
	var buf bytes.Buffer

	for _, char := range str {
		for _, sp := range specials {
			if sp == char {
				buf.WriteRune('\\')

				continue
			}
		}

		buf.WriteRune(char)
	}

	return buf.String()
}

func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	truncated := ""
	count := 0

	for _, char := range str {
		truncated += string(char)
		count++

		if count >= length {
			break
		}
	}

	return truncated
}
