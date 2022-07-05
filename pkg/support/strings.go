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
