package rules

// TODO: Create more flexible approach to create correct execution order
var (
	BeforeScripts byte = 1
	AfterScripts  byte = 2
)

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

func (rule BaseRule) GetPosition() byte {
	return BeforeScripts
}
