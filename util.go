package goslack

import "strings"

// StringPop splits the string on the first space and return both parts trimmed of whitespace
func StringPop(m string) (first string, rest string) {
	parts := strings.SplitN(m, " ", 2)

	if len(parts) < 1 {
		return
	}
	first = parts[0]

	rest = ""
	if len(parts) > 1 {
		rest = strings.Trim(parts[1], " \t\r\n")
	}

	return
}
