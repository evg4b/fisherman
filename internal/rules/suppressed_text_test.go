package rules_test

import (
	"context"
	"errors"
	"fisherman/internal"
	"fisherman/pkg/vcs"
	"fisherman/testing/mocks"
	"io"
	"strings"
	"testing"

	. "fisherman/internal/rules"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuppressedText_Check(t *testing.T) {
	t.Run("checked suppression", func(t *testing.T) {
		tests := []struct {
			name        string
			substrings  []string
			repo        internal.Repository
			expectedErr []string
		}{
			{
				name: "no changes",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.Return(map[string]vcs.Changes{}, nil),
				substrings: []string{"suppressed text"},
			},
			{
				name: "suppressed text not found",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {{Status: vcs.Added, Change: "this is valid text"}},
						"README.md":    {{Status: vcs.Added, Change: "hello word"}},
					}, nil),
				substrings: []string{"suppressed text"},
			},
			{
				name: "suppressed text is deleted",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {{Status: vcs.Deleted, Change: "this is suppressed text"}},
						"README.md":    {{Status: vcs.Added, Change: "hello word"}},
					}, nil),
				substrings: []string{"suppressed text"},
			},
			{
				name: "suppressed text is deleted",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {{Status: vcs.Deleted, Change: "this is suppressed text"}},
						"README.md":    {{Status: vcs.Added, Change: "hello word"}},
					}, nil),
				substrings: []string{"suppressed text"},
			},
			{
				name: "suppressed text founded",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {{Status: vcs.Added, Change: "this is suppressed text"}},
						"README.md":    {{Status: vcs.Added, Change: "hello word"}},
					}, nil),
				substrings:  []string{"suppressed text"},
				expectedErr: []string{"file 'test/file.go' should not contains 'suppressed text'"},
			},
			{
				name: "suppressed multiple text founded",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {
							{Status: vcs.Added, Change: "this is suppressed text"},
							{Status: vcs.Added, Change: "this is second suppressed text line"},
						},
						"README.md": {{Status: vcs.Added, Change: "hello word"}},
					}, nil),
				substrings:  []string{"suppressed text"},
				expectedErr: []string{"file 'test/file.go' should not contains 'suppressed text'"},
			},
			{
				name: "multiple files with suppressed text founded",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {{Status: vcs.Added, Change: "this is suppressed text"}},
						"README.md":    {{Status: vcs.Added, Change: "this is second suppressed text in other file"}},
					}, nil),
				substrings: []string{"suppressed text"},
				expectedErr: []string{
					"file 'test/file.go' should not contains 'suppressed text'",
					"file 'README.md' should not contains 'suppressed text'",
				},
			},
			{
				name: "multiple suppressed string in one line",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {{Status: vcs.Added, Change: "this is suppressed text"}},
						"README.md":    {{Status: vcs.Deleted, Change: "this is second suppressed text in other file"}},
					}, nil),
				substrings: []string{"suppressed", "text"},
				expectedErr: []string{
					"file 'test/file.go' should not contains 'suppressed'",
					"file 'test/file.go' should not contains 'text'",
				},
			},
			{
				name: "multiple suppressed string in difference lines",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {
							{Status: vcs.Added, Change: "this is suppressed line"},
							{Status: vcs.Added, Change: "this is not allowed text"},
						},
						"README.md": {
							{Status: vcs.Added, Change: "this is suppressed line"},
							{Status: vcs.Added, Change: "this is not allowed text"},
						},
					}, nil),
				substrings: []string{"suppressed", "text"},
				expectedErr: []string{
					"file 'test/file.go' should not contains 'suppressed'",
					"file 'test/file.go' should not contains 'text'",
					"file 'README.md' should not contains 'suppressed'",
					"file 'README.md' should not contains 'text'",
				},
			},
			{
				name: "internal error",
				repo: mocks.NewRepositoryMock(t).
					GetIndexChangesMock.
					Return(map[string]vcs.Changes{}, errors.New("test error")),
				substrings:  []string{"suppressed", "text"},
				expectedErr: []string{"test error"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				rule := makeRule(
					&SuppressedText{
						BaseRule:   BaseRule{Type: SuppressedTextType},
						Substrings: tt.substrings,
					},
					WithRepository(tt.repo),
				)

				err := rule.Check(context.TODO(), io.Discard)

				assertFlatMultiError(t, tt.expectedErr, err)
			})
		}
	})

	t.Run("checked suppression with excluded files", func(t *testing.T) {
		tests := []struct {
			name        string
			exclude     []string
			expectedErr []string
		}{
			{
				name:    "suppressed files not excluded",
				exclude: []string{"other-file.go"},
				expectedErr: []string{
					"file 'test/file.go' should not contains 'suppressed text'",
					"file 'README.md' should not contains 'suppressed text'",
				},
			},
			{
				name:        "suppressed single files excluded with glob",
				exclude:     []string{"*.md"},
				expectedErr: []string{"file 'test/file.go' should not contains 'suppressed text'"},
			},
			{
				name:        "suppressed single files not excluded with path",
				exclude:     []string{"test/file.go"},
				expectedErr: []string{"file 'README.md' should not contains 'suppressed text'"},
			},
			{
				name:        "winvalid glob pattern",
				exclude:     []string{"some/[*"},
				expectedErr: []string{"syntax error in pattern"},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repo := mocks.NewRepositoryMock(t).GetIndexChangesMock.
					Return(map[string]vcs.Changes{
						"test/file.go": {
							{Status: vcs.Added, Change: "this is suppressed text"},
							{Status: vcs.Added, Change: "this is allowed text"},
						},
						"README.md": {
							{Status: vcs.Added, Change: "this is suppressed text"},
							{Status: vcs.Added, Change: "this is allowed text"},
						},
					}, nil)

				rule := makeRule(
					&SuppressedText{
						BaseRule:      BaseRule{Type: SuppressedTextType},
						Substrings:    []string{"suppressed text"},
						ExcludedGlobs: tt.exclude,
					},
					WithRepository(repo),
				)

				err := rule.Check(context.TODO(), io.Discard)

				assertFlatMultiError(t, tt.expectedErr, err)
			})
		}
	})
}

func TestSuppressedText_GetPosition(t *testing.T) {
	rule := SuppressedText{
		BaseRule:   BaseRule{Type: SuppressedTextType},
		Substrings: []string{"suppressed text"},
	}

	actual := rule.GetPosition()

	assert.Equal(t, actual, PostScripts)
}

func TestSuppressedText_Compile(t *testing.T) {
	rule := SuppressedText{
		BaseRule:   BaseRule{Type: ShellScriptType},
		Substrings: []string{"DEMO {{var1}}"},
		ExcludedGlobs: []string{
			"Glob1 {{var1}}",
			"Glob2 {{var1}}",
		},
	}

	rule.Compile(map[string]any{"var1": "VALUE"})

	assert.Equal(t, SuppressedText{
		BaseRule:   BaseRule{Type: ShellScriptType},
		Substrings: []string{"DEMO VALUE"},
		ExcludedGlobs: []string{
			"Glob1 VALUE",
			"Glob2 VALUE",
		},
	}, rule)
}

func assertFlatMultiError(t *testing.T, expected []string, actual error) {
	t.Helper()

	if len(expected) > 0 {
		message := actual.Error()
		parts := strings.Split(message, "\n")
		assert.Equal(t, len(parts), len(expected))
		for _, expectedLine := range expected {
			assert.Contains(t, parts, expectedLine)
		}
	} else {
		require.NoError(t, actual)
	}
}
