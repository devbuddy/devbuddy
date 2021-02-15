package context

import "regexp"

// const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
const ansi = `\x1B\[\d+(;\d+){0,2}m`

var re = regexp.MustCompile(ansi)

func StripAnsi(str string) string {
	return re.ReplaceAllString(str, "")
}

func StripAnsiSlice(slice []string) []string {
	var res []string
	for _, s := range slice {
		res = append(res, StripAnsi(s))
	}
	return res
}
