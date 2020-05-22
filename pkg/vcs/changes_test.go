package vcs_test

import (
	"testing"

	. "github.com/evg4b/fisherman/pkg/vcs"

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
			Change{Status: Added, Change: "test"},
			Change{Status: Added, Change: "data"},
			Change{Status: Added, Change: "diff"},
		},
		added: Changes{
			Change{Status: Added, Change: "test"},
			Change{Status: Added, Change: "data"},
			Change{Status: Added, Change: "diff"},
		},
		deleted:    Changes{},
		unmodified: Changes{},
	},
	{
		name: "collection 2",
		changes: Changes{
			Change{Status: Deleted, Change: "test"},
			Change{Status: Deleted, Change: "data"},
			Change{Status: Deleted, Change: "diff"},
		},
		deleted: Changes{
			Change{Status: Deleted, Change: "test"},
			Change{Status: Deleted, Change: "data"},
			Change{Status: Deleted, Change: "diff"},
		},
		added:      Changes{},
		unmodified: Changes{},
	},
	{
		name: "collection 3",
		changes: Changes{
			Change{Status: Unmodified, Change: "test"},
			Change{Status: Unmodified, Change: "data"},
			Change{Status: Unmodified, Change: "diff"},
		},
		unmodified: Changes{
			Change{Status: Unmodified, Change: "test"},
			Change{Status: Unmodified, Change: "data"},
			Change{Status: Unmodified, Change: "diff"},
		},
		added:   Changes{},
		deleted: Changes{},
	},
	{
		name: "collection 4",
		changes: Changes{
			Change{Status: Unmodified, Change: "test"},
			Change{Status: Added, Change: "data"},
			Change{Status: Deleted, Change: "diff"},
		},
		unmodified: Changes{
			Change{Status: Unmodified, Change: "test"},
		},
		added: Changes{
			Change{Status: Added, Change: "data"},
		},
		deleted: Changes{
			Change{Status: Deleted, Change: "diff"},
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
	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			actual := testCase.changes.Deleted()

			assert.Equal(t, testCase.deleted, actual)
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
