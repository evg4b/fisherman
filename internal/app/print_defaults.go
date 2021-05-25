package app

import (
	"fisherman/internal/constants"
	"fisherman/internal/utils"
	"fisherman/pkg/log"
	"fmt"
	"strings"
)

const preffix = "  "
const baseCommandLen = 8

func (r *FishermanApp) PrintDefaults() {
	utils.PrintGraphics(log.InfoOutput, constants.Logo, map[string]interface{}{
		constants.FishermanVersionVariable: constants.Version,
	})
	fmt.Fprintln(log.InfoOutput, "                 Small git hook management tool for developer.")
	fmt.Fprintln(log.InfoOutput, "")
	fmt.Fprintln(log.InfoOutput, preffix, "Commands :")
	for _, command := range r.commands {
		fmt.Fprintln(
			log.InfoOutput,
			preffix,
			preffix,
			command.Name(),
			strings.Repeat(" ", baseCommandLen-len(command.Name())),
			command.Description())
	}
	fmt.Fprintln(log.InfoOutput, "")
}
