package commands

type CliCommand interface {
	Init([]string) error
	Run(ctx Context) error
	Name() string
}
