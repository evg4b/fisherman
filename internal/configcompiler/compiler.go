package configcompiler

import (
	"fisherman/configuration"
	"fisherman/infrastructure"
)

type Compiler = func(config CompilableConfig)

type CompilableConfig interface {
	Compile(configuration.Variables)
	GetVariablesConfig() configuration.VariablesConfig
}

func NewCompiler(repository infrastructure.Repository, globalVars configuration.Variables, cwd string) Compiler {
	extractor := NewConfigExtractor(repository, globalVars, cwd)

	return func(config CompilableConfig) {
		variables, err := extractor.Variables(config.GetVariablesConfig())
		if err != nil {
			panic(err)
		}

		config.Compile(variables)
	}
}
