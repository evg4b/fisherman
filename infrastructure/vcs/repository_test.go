package vcs

import (
	"fisherman/infrastructure"
	"reflect"
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestNewGitRepository(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *GitRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGitRepository(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGitRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepository_GetCurrentBranch(t *testing.T) {
	tests := []struct {
		name    string
		r       *GitRepository
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetCurrentBranch()
			if (err != nil) != tt.wantErr {
				t.Errorf("GitRepository.GetCurrentBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GitRepository.GetCurrentBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepository_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		r       *GitRepository
		want    infrastructure.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetUser()
			if (err != nil) != tt.wantErr {
				t.Errorf("GitRepository.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitRepository.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepository_repo(t *testing.T) {
	tests := []struct {
		name string
		r    *GitRepository
		want *git.Repository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.repo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitRepository.repo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepository_AddGlob(t *testing.T) {
	type args struct {
		glob string
	}
	tests := []struct {
		name    string
		r       *GitRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.AddGlob(tt.args.glob); (err != nil) != tt.wantErr {
				t.Errorf("GitRepository.AddGlob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGitRepository_RemoveGlob(t *testing.T) {
	type args struct {
		glob string
	}
	tests := []struct {
		name    string
		r       *GitRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.RemoveGlob(tt.args.glob); (err != nil) != tt.wantErr {
				t.Errorf("GitRepository.RemoveGlob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
