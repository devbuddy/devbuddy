package colorfmt

import "github.com/dolab/colorize"

const tagResetColor = "\x1b[0m"

type tagProcessor func(tagContent string) string

type colorizeTagProcessor struct {
	specialColors map[string]string
}

func newColorizeTagProcessor() *colorizeTagProcessor {
	return &colorizeTagProcessor{
		specialColors: map[string]string{
			"link": "green+ubh",
		},
	}
}

func (c *colorizeTagProcessor) process(text string) string {
	if text == "reset" {
		return tagResetColor
	}
	if translated, ok := c.specialColors[text]; ok {
		text = translated
	}
	colorDraw, _ := colorize.New(text).Colour()
	return colorDraw
}

func ignoreTagProcessor(_ string) string {
	return ""
}
