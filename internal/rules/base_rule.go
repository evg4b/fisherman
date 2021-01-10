package rules

type BaseRule struct {
	Type      string `yaml:"type,omitempty"`
	Condition string `yaml:"condition,omitempty"`
}

func (rule BaseRule) GetType() string {
	return rule.Type
}

func (rule BaseRule) GetContition() string {
	return rule.Condition
}
