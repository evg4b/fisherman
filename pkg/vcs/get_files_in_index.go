package vcs

func (r *GitRepository) GetFilesInIndex() ([]string, error) {
	if err := r.init(); err != nil {
		return nil, err
	}

	worktree, err := r.repo.Worktree()
	if err != nil {
		return nil, err
	}

	statusIndex, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	files := []string{}
	for key, status := range statusIndex {
		if isModified(status.Staging) {
			files = append(files, key)
		}
	}

	return files, nil
}
