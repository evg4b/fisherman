package utils

import (
	"fisherman/pkg/guards"
	"io"

	"github.com/valyala/fasttemplate"
)

const (
	startTag = "{{"
	endTag   = "}}"
)

func PrintGraphics(wr io.Writer, content string, data map[string]any) {
	tpl := makeTemplate(content)
	_, err := tpl.Execute(wr, data)
	guards.NoError(err)
}

func FillTemplate(src *string, data map[string]any) {
	(*src) = makeTemplate(*src).ExecuteString(data)
}

func FillTemplatesArray(src []string, data map[string]any) {
	for index, srcItem := range src {
		src[index] = makeTemplate(srcItem).ExecuteString(data)
	}
}

func FillTemplatesMap(src map[string]string, data map[string]any) {
	for key, srcItem := range src {
		src[key] = makeTemplate(srcItem).ExecuteString(data)
	}
}

func makeTemplate(content string) *fasttemplate.Template {
	return fasttemplate.New(content, startTag, endTag)
}
