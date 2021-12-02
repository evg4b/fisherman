package vcs_test

import (
	. "fisherman/pkg/vcs"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name       string
	changes    Changes
	added      Changes
	deleted    Changes
	unmodified Changes
}{
	{
		name:       "empty collection",
		changes:    Changes{},
		added:      Changes{},
		deleted:    Changes{},
		unmodified: Changes{},
	},
	{
		name: "collection 1",
		changes: Changes{
			Change{Added, "test"},
			Change{Added, "data"},
			Change{Added, "diff"},
		},
		added: Changes{
			Change{Added, "test"},
			Change{Added, "data"},
			Change{Added, "diff"},
		},
		deleted:    Changes{},
		unmodified: Changes{},
	},
	{
		name: "collection 2",
		changes: Changes{
			Change{Deleted, "test"},
			Change{Deleted, "data"},
			Change{Deleted, "diff"},
		},
		deleted: Changes{
			Change{Deleted, "test"},
			Change{Deleted, "data"},
			Change{Deleted, "diff"},
		},
		added:      Changes{},
		unmodified: Changes{},
	},
	{
		name: "collection 3",
		changes: Changes{
			Change{Unmodified, "test"},
			Change{Unmodified, "data"},
			Change{Unmodified, "diff"},
		},
		unmodified: Changes{
			Change{Unmodified, "test"},
			Change{Unmodified, "data"},
			Change{Unmodified, "diff"},
		},
		added:   Changes{},
		deleted: Changes{},
	},
	{
		name: "collection 4",
		changes: Changes{
			Change{Unmodified, "test"},
			Change{Added, "data"},
			Change{Deleted, "diff"},
		},
		unmodified: Changes{
			Change{Unmodified, "test"},
		},
		added: Changes{
			Change{Added, "data"},
		},
		deleted: Changes{
			Change{Deleted, "diff"},
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
