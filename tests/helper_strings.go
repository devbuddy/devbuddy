package integration

import "regexp"

const stripANSI = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var reStripANSI = regexp.MustCompile(stripANSI)

func StripANSI(str string) string {
	return reStripANSI.ReplaceAllString(str, "")
}
