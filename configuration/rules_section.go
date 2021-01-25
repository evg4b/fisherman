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

	switch typeString.(string) {
	case rules.SuppressCommitType:
		var rule rules.SuppressCommitFiles
		err := decode(rawRule, &rule)
		if err != nil {
			return nil, err
		}

		return rule, nil

	case rules.CommitMessageType:
		var rule rules.CommitMessage
		err := decode(rawRule, &rule)
		if err != nil {
			return nil, err
		}

		return rule, nil

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
