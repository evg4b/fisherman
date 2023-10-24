package app

import (
	"fmt"
	"github.com/evg4b/fisherman/internal/constants"
	"github.com/evg4b/fisherman/internal/utils"
	"github.com/evg4b/fisherman/pkg/log"
	"strings"
)

const (
	prefix         = "  "
	baseCommandLen = 8
)

// PrintDefaults prints custom information output.
func (r *FishermanApp) PrintDefaults() {
	utils.PrintGraphics(log.InfoOutput, constants.Logo, map[string]any{
		constants.FishermanVersionVariable: constants.Version,
	})

	fmt.Fprintln(log.InfoOutput, strings.Repeat(prefix, 8), "Small git hook management tool for developer.") // nolint mnd
	fmt.Fprintln(log.InfoOutput, "")
	fmt.Fprintln(log.InfoOutput, prefix, "Commands:")

	for _, command := range r.commands {
		fmt.Fprintln(
			log.InfoOutput,
			strings.Repeat(prefix, 2), // nolint mnd
			command.Name(),
			strings.Repeat(" ", baseCommandLen-len(command.Name())),
			command.Description())
	}

	fmt.Fprintln(log.InfoOutput, "")
	fmt.Fprintln(log.InfoOutput, prefix, "Configuration docs:", constants.ConfigurationDocksURL)
	fmt.Fprintln(log.InfoOutput, "")
}
