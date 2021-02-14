package utils

import (
	"io"

	"github.com/valyala/fasttemplate"
)

const startTag = "{{"
const endTag = "}}"

func PrintGraphics(wr io.Writer, content string, data map[string]interface{}) {
	tpl := fasttemplate.New(content, startTag, endTag)
	_, err := tpl.Execute(wr, data)
	HandleCriticalError(err)
}

func FillTemplate(src *string, data map[string]interface{}) {
	(*src) = fasttemplate.New(*src, startTag, endTag).ExecuteString(data)
}

func FillTemplatesArray(src []string, data map[string]interface{}) {
	for index, srcItem := range src {
		src[index] = fasttemplate.New(srcItem, startTag, endTag).ExecuteString(data)
	}
}

func FillTemplatesMap(src map[string]string, data map[string]interface{}) {
	for key, srcItem := range src {
		src[key] = fasttemplate.New(srcItem, startTag, endTag).ExecuteString(data)
	}
}
