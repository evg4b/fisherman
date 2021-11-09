package app

import (
	"fisherman/internal/constants"
	"fisherman/internal/utils"
	"fisherman/pkg/log"
	"fmt"
	"strings"
)

const (
	preffix        = "  "
	baseCommandLen = 8
)

// PrintDefaults prints custom information output.
func (r *FishermanApp) PrintDefaults() {
	utils.PrintGraphics(log.InfoOutput, constants.Logo, map[string]interface{}{
		constants.FishermanVersionVariable: constants.Version,
	})

	fmt.Fprintln(log.InfoOutput, strings.Repeat(preffix, 8), "Small git hook management tool for developer.") // nolint mnd
	fmt.Fprintln(log.InfoOutput, "")
	fmt.Fprintln(log.InfoOutput, preffix, "Commands:")

	for _, command := range r.commands {
		fmt.Fprintln(
			log.InfoOutput,
			strings.Repeat(preffix, 2), // nolint mnd
			command.Name(),
			strings.Repeat(" ", baseCommandLen-len(command.Name())),
			command.Description())
	}

	fmt.Fprintln(log.InfoOutput, "")
	fmt.Fprintln(log.InfoOutput, preffix, "Configuration docs:", constants.ConfigurationDocksURL)
	fmt.Fprintln(log.InfoOutput, "")
}
