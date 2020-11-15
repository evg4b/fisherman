package utils

import (
	"io"

	"github.com/valyala/fasttemplate"
)

func PrintGraphics(wr io.Writer, content string, data map[string]interface{}) {
	tpl := fasttemplate.New(content, "{{", "}}")
	_, err := tpl.Execute(wr, data)
	HandleCriticalError(err)
}

func FillTemplate(src *string, data map[string]interface{}) {
	*src = fasttemplate.
		New(*src, "{{", "}}").
		ExecuteString(data)
}
