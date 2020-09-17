package utils

import (
	"html/template"
	"io"
)

// PrintGraphics prints fill template data from map or object and put this content in io.Writer.
func PrintGraphics(wr io.Writer, content string, data interface{}) {
	tpl, err := template.New(content).Parse(content)
	if err != nil {
		panic(err)
	}

	if err := tpl.Execute(wr, data); err != nil {
		panic(err)
	}
}
