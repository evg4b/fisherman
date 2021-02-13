package utils

import (
	"io"

	"github.com/imdario/mergo"
	"github.com/valyala/fasttemplate"
)

const startTag = "{{"
const endTag = "}}"

func PrintGraphics(wr io.Writer, content string, data map[string]interface{}) {
	tpl := fasttemplate.New(content, startTag, endTag)
	_, err := tpl.Execute(wr, data)
	HandleCriticalError(err)
}

func FillTemplate(src string, data ...map[string]interface{}) string {
	dest := map[string]interface{}{}
	for _, src := range data {
		err := mergo.MergeWithOverwrite(&dest, src)
		if err != nil {
			panic(err)
		}
	}

	return fasttemplate.New(src, startTag, endTag).ExecuteString(dest)
}
