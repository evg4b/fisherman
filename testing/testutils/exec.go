package testutils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type (
	cmdBuilder = func(context.Context, string, ...string) *exec.Cmd
	envBuilder = func([]string) []string
)

const envConstant = "GO_HELPER_PROCESS"

// ConfigureFakeExec configures exec.Cmd builder mock to call passed helper test
// (helper test should be located in the package with test).
func ConfigureFakeExec(helperTest string) (cmdBuilder, envBuilder) {
	cmdBuilder := func(_ context.Context, program string, args ...string) *exec.Cmd {
		testArgs := []string{"-test.run=" + helperTest, "-test.v", "--", program}
		testArgs = append(testArgs, args...)

		return exec.Command(os.Args[0], testArgs...) // nolint gosec
	}

	envBuilder := func(env []string) []string {
		if env == nil {
			env = []string{}
		}

		return append(env, envConstant+"=1")
	}

	return cmdBuilder, envBuilder
}

// ExecTestHandler configure helper test to check executed command.
// Key is concatanaion program with arguments e.g. 'go test ./..'.
func ExecTestHandler(t *testing.T, cases map[string]func()) {
	t.Helper()

	if os.Getenv(envConstant) != "1" {
		t.Skip("helper test")
	}

	testCase := strings.Join(os.Args[4:], " ")

	caseFunction, ok := cases[testCase]
	if !ok {
		panic(fmt.Sprintf("test '%s' case not specified", testCase))
	}

	caseFunction()
}
