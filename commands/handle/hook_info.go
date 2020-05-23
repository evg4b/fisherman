package handle

import (
	"fisherman/constants"
	"io"
	"text/template"
)

type HookInfo struct {
	Hook             string
	GlobalConfigPath string
	RepoConfigPath   string
	LocalConfigPath  string
	Version          string
}

func printHookHeader(info *HookInfo, wr io.Writer) {
	tmpl, err := template.New("hook-logo").Parse(constants.HookHeader)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(wr, info)
	if err != nil {
		panic(err)
	}
}
