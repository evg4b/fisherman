package shell_test

const windowsOS = "windows"

var scriptTests = []struct {
	name        string
	script      string
	expectedRrr string
}{
	{
		name:   "exec successfully",
		script: "echo 'test'",
	},
	{
		name:        "exec with error",
		script:      "exit 33",
		expectedRrr: "exit status 33",
	},
}

var envTests = []struct {
	name     string
	env      []string
	expected []string
}{
	{
		name:     "empty slice",
		env:      []string{},
		expected: []string{},
	},
	{
		name:     "additional arguments",
		env:      []string{"VAR1=1", "VAR2=2"},
		expected: []string{"VAR1=1", "VAR2=2"},
	},
}
