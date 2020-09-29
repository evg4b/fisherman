package runner

import (
	"fisherman/constants"
	"fisherman/utils"

	"github.com/imdario/mergo"
)

func (r *Runner) buildVariables() (map[string]interface{}, error) {
	gitUser, err := r.repository.GetUser()
	if err != nil {
		return nil, err
	}

	variables := map[string]interface{}{
		"FishermanVersion": constants.Version,
		"CWD":              r.app.Cwd,
		"UserName":         gitUser.UserName,
		"Email":            gitUser.Email,
	}

	err = mergo.Map(&variables, r.config.GlobalVariables)
	utils.HandleCriticalError(err)

	return variables, nil
}
