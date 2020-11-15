package configcompiler

import (
	"fisherman/config/hooks"
	"fisherman/infrastructure"
)

type Compiler = func(config CompilableConfig)

type CompilableConfig interface {
	Compile(map[string]interface{})
	GetVarsSection() hooks.Variables
	HasVars() bool
}

func NewCompiler(repository infrastructure.Repository, globalVariables map[string]interface{}, cwd string) Compiler {
	extractor := NewConfigExtractor(repository, globalVariables, cwd)

	return func(config CompilableConfig) {
		variables, err := extractor.Variables(config.GetVarsSection())
		if err != nil {
			panic(err)
		}

		config.Compile(variables)
	}
}
