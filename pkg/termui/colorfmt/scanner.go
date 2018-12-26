package colorfmt

func scan(text string, openingChar, closingChar rune, processor tagProcessor) string {
	var in bool
	var tag string
	var output string

	for _, char := range text {
		if in {
			if char == '{' && tag == "" { // "{{"
				in = false
				output += string(char)
			} else if char == closingChar {
				output += processor(tag)
				in = false
				tag = ""
			} else {
				tag += string(char)
			}
		} else {
			if char == openingChar {
				in = true
			} else {
				output += string(char)
			}
		}
	}

	return output
}
