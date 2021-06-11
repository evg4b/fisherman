package configuration

import (
	"fisherman/internal/rules"
	"fisherman/internal/utils"
	"reflect"

	"github.com/go-errors/errors"
	"gopkg.in/yaml.v3"
)

type ruleDef struct {
	Type string `yaml:"type"`
	Rule Rule
}

var definedTypes = map[string]reflect.Type{
	rules.SuppressCommitFilesType: reflect.TypeOf(rules.SuppressCommitFiles{}),
	rules.CommitMessageType:       reflect.TypeOf(rules.CommitMessage{}),
	rules.PrepareMessageType:      reflect.TypeOf(rules.PrepareMessage{}),
	rules.ShellScriptType:         reflect.TypeOf(rules.ShellScript{}),
	rules.AddToIndexType:          reflect.TypeOf(rules.AddToIndex{}),
}

func (def *ruleDef) UnmarshalYAML(value *yaml.Node) error {
	type plain ruleDef
	err := value.Decode((*plain)(def))
	if err != nil {
		return err
	}

	if utils.IsEmpty(def.Type) {
		return errors.Errorf("required property 'type' not defined")
	}

	reflectType, ok := definedTypes[def.Type]
	if !ok {
		return errors.Errorf("type %s is not supported", def.Type)
	}

	def.Rule = reflect.New(reflectType).Interface().(Rule)

	return value.Decode(def.Rule)
}
