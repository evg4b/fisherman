package shell

import "time"

type Script struct {
	commands []string
	env      map[string]string
	dir      string
	duration time.Duration
}

// TODO: Add options pattern to configure script.
func NewScript(commands []string) *Script {
	return &Script{
		commands: commands,
		env:      map[string]string{},
	}
}

func (s *Script) GetCommands() []string {
	return s.commands
}

func (s *Script) SetEnvironmentVariables(env map[string]string) *Script {
	s.env = env

	return s
}

func (s *Script) GetEnvironmentVariables() map[string]string {
	return s.env
}

func (s *Script) GetDirectory() string {
	return s.dir
}

func (s *Script) SetDirectory(dir string) *Script {
	s.dir = dir

	return s
}

func (s *Script) GetDuration() time.Duration {
	return s.duration
}
