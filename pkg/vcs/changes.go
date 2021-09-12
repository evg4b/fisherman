package vcs

type ChangeCode byte

const (
	Unmodified ChangeCode = '='
	Added      ChangeCode = '+'
	Deleted    ChangeCode = '-'
)

type Change struct {
	Status ChangeCode
	Change string
}

type Changes []Change

func (changes Changes) Added() Changes {
	return changes.filter(Added)
}

func (changes Changes) Deleted() Changes {
	return changes.filter(Deleted)
}

func (changes Changes) filter(status ChangeCode) (filtered Changes) {
	for _, change := range changes {
		if change.Status == status {
			filtered = append(filtered, change)
		}
	}

	return filtered
}
