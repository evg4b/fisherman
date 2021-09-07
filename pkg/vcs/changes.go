package vcs

type ChangeCode byte

const (
	Unmodified ChangeCode = '='
	Added      ChangeCode = '+'
	Deleted    ChangeCode = '-'
)

type Change struct {
	Status ChangeCode
	Change string
}

type Changes []Change
