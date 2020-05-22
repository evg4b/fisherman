package commands

type CliCommand interface {
	Init([]string) error
	Run() error
	Name() string
}
