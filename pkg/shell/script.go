package shell

type Script struct {
	commands []string
	env      map[string]string
	dir      string
}

func NewScript() *Script {
	return &Script{
		commands: []string{},
		env:      map[string]string{},
	}
}

func (s *Script) SetCommands(commands []string) *Script {
	s.commands = commands

	return s
}

func (s *Script) SetEnvironmentVariables(env map[string]string) *Script {
	s.env = env

	return s
}

func (s *Script) SetDirectory(dir string) *Script {
	s.dir = dir

	return s
}
