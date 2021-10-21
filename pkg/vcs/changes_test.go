package vcs_test

import (
	"fisherman/pkg/vcs"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name       string
	changes    vcs.Changes
	added      vcs.Changes
	deleted    vcs.Changes
	unmodified vcs.Changes
}{
	{
		name:       "empty collection",
		changes:    vcs.Changes{},
		added:      vcs.Changes{},
		deleted:    vcs.Changes{},
		unmodified: vcs.Changes{},
	},
	{
		name: "collection 1",
		changes: vcs.Changes{
			vcs.Change{vcs.Added, "test"},
			vcs.Change{vcs.Added, "data"},
			vcs.Change{vcs.Added, "diff"},
		},
		added: vcs.Changes{
			vcs.Change{vcs.Added, "test"},
			vcs.Change{vcs.Added, "data"},
			vcs.Change{vcs.Added, "diff"},
		},
		deleted:    vcs.Changes{},
		unmodified: vcs.Changes{},
	},
	{
		name: "collection 2",
		changes: vcs.Changes{
			vcs.Change{vcs.Deleted, "test"},
			vcs.Change{vcs.Deleted, "data"},
			vcs.Change{vcs.Deleted, "diff"},
		},
		deleted: vcs.Changes{
			vcs.Change{vcs.Deleted, "test"},
			vcs.Change{vcs.Deleted, "data"},
			vcs.Change{vcs.Deleted, "diff"},
		},
		added:      vcs.Changes{},
		unmodified: vcs.Changes{},
	},
	{
		name: "collection 3",
		changes: vcs.Changes{
			vcs.Change{vcs.Unmodified, "test"},
			vcs.Change{vcs.Unmodified, "data"},
			vcs.Change{vcs.Unmodified, "diff"},
		},
		unmodified: vcs.Changes{
			vcs.Change{vcs.Unmodified, "test"},
			vcs.Change{vcs.Unmodified, "data"},
			vcs.Change{vcs.Unmodified, "diff"},
		},
		added:   vcs.Changes{},
		deleted: vcs.Changes{},
	},
	{
		name: "collection 4",
		changes: vcs.Changes{
			vcs.Change{vcs.Unmodified, "test"},
			vcs.Change{vcs.Added, "data"},
			vcs.Change{vcs.Deleted, "diff"},
		},
		unmodified: vcs.Changes{
			vcs.Change{vcs.Unmodified, "test"},
		},
		added: vcs.Changes{
			vcs.Change{vcs.Added, "data"},
		},
		deleted: vcs.Changes{
			vcs.Change{vcs.Deleted, "diff"},
		},
	},
}

func TestChanges_Added(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.changes.Added()

			assert.Equal(t, tt.added, actual)
		})
	}
}

func TestChanges_Deleted(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.changes.Deleted()

			assert.Equal(t, tt.deleted, actual)
		})
	}
}

func TestChanges_Unmodified(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.changes.Unmodified()

			assert.Equal(t, tt.unmodified, actual)
		})
	}
}
