package utils

import "path/filepath"

func MatchToGlobs(patterns []string, file string) (matched bool, err error) {
	if len(patterns) > 0 {
		for _, glob := range patterns {
			matched, err := filepath.Match(glob, file)
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
	}

	return false, nil
}
