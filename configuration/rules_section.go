package configuration

import (
	"errors"
	"fisherman/internal"
	"fisherman/internal/rules"
	"fmt"
	"io"

	"github.com/mitchellh/mapstructure"
)

const rulesKey = "rules"
const typeKey = "type"

// TODO: Add new method in Rule interface to Decode rule from map[string]interface{} and
// try implement comman realization in base rule structure.
type Rule interface {
	GetType() string
	GetContition() string
	GetPosition() byte
	Check(io.Writer, internal.ExecutionContext) error
}

type RulesSection struct {
	Rules []Rule
}

func (config *RulesSection) UnmarshalYAML(unmarshal func(interface{}) error) error {
	config.Rules = []Rule{}
	var rawSection map[string]interface{}

	err := unmarshal(&rawSection)
	if err != nil {
		return err
	}

	rulesSection, ok := rawSection[rulesKey]
	if !ok {
		return nil
	}

	rawRules, ok := rulesSection.([]interface{})
	if !ok {
		return errors.New("unknown rules markup")
	}

	for index, rawRule := range rawRules {
		rule, err := unmarshalRule(rawRule)
		if err != nil {
			return fmt.Errorf("error for rule at index %d: %w", index, err)
		}

		config.Rules = append(config.Rules, rule)
	}

	return nil
}

func unmarshalRule(rawRule interface{}) (Rule, error) {
	typeString, ok := rawRule.(map[string]interface{})[typeKey]
	if !ok {
		return nil, fmt.Errorf("required property '%s' not defined", typeKey)
	}

	rule, err := selectRule(typeString.(string))
	if err != nil {
		return nil, err
	}

	err = decode(rawRule, rule)

	return rule, err
}

func selectRule(typeName string) (Rule, error) {
	switch typeName {
	case rules.SuppressCommitFilesType:
		return &rules.SuppressCommitFiles{}, nil
	case rules.CommitMessageType:
		return &rules.CommitMessage{}, nil
	case rules.PrepareMessageType:
		return &rules.PrepareMessage{}, nil
	case rules.ShellScriptType:
		return &rules.ShellScript{}, nil
	default:
		return nil, errors.New("unknown rule type")
	}
}

func decode(input interface{}, output interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      output,
		ErrorUnused: true,
	})

	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
