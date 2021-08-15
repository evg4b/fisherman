package vcs

import "github.com/go-git/go-git/v5"

type Change struct {
	Status git.StatusCode
	Change string
}

type Changes []Change
