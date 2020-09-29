package utils

import (
	"io"

	"github.com/valyala/fasttemplate"
)

// PrintGraphics prints fill template data from map or object and put this content in io.Writer.
func PrintGraphics(wr io.Writer, content string, data map[string]interface{}) {
	tpl := fasttemplate.New(content, "{{", "}}")
	_, err := tpl.Execute(wr, data)
	HandleCriticalError(err)
}
