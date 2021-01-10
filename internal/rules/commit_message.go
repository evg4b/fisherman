package rules

const CommitMessageType = "commit-message"

type CommitMessage struct {
	BaseRule `mapstructure:",squash"`
	Prefix   string `mapstructure:"prefix"`
	Suffix   string `mapstructure:"suffix"`
	Regexp   string `mapstructure:"regexp"`
}
