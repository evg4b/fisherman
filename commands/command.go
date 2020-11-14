package commands

type CliCommand interface {
	Init(args []string) error
	Run() error
	Name() string
	Description() string
}
