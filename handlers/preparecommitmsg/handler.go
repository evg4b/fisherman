package preparecommitmsg

import (
	"fisherman/commands/context"
	"fisherman/infrastructure/io"
	"regexp"
	"strings"
)

// Handler is structure for storage information about prepare-commit-msg hook handler
type Handler struct {
	fileAccessor io.FileAccessor
}

// NewHandler is constructor for prepare-commit-msg hook handler
func NewHandler(fileAccessor io.FileAccessor) *Handler {
	return &Handler{fileAccessor}
}

// Execute is a execute function for prepare-commit-msg hook
func (h *Handler) Execute(ctx context.Context, args []string) {
	config := ctx.GetConfiguration()
	info, err := ctx.GetGitInfo()
	if err != nil {
		panic(err)
	}

	c := config.Hooks.PrepareCommitMsgHook
	if c != nil {
		message, isPresented := getPreparedMessage(c.Message, c.BranchRegExp, info.CurrentBranch)
		if isPresented {
			err = h.fileAccessor.WriteFile(args[0], message)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getPreparedMessage(message, regexpString, branch string) (string, bool) {
	if !isEmpty(message) {
		if !isEmpty(regexpString) {
			return regexp.MustCompile(regexpString).
				ReplaceAllString(branch, message), true
		}

		return message, true
	}
	return "", false
}

func isEmpty(data string) bool {
	return len(strings.TrimSpace(data)) == 0
}
