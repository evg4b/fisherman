package vcs

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage"
)

// nolint: interfacebloat
type GoGitRepository interface {
	BlobObject(h plumbing.Hash) (*object.Blob, error)
	BlobObjects() (*object.BlobIter, error)
	Branch(name string) (*config.Branch, error)
	Branches() (storer.ReferenceIter, error)
	CommitObject(h plumbing.Hash) (*object.Commit, error)
	CommitObjects() (object.CommitIter, error)
	Config() (*config.Config, error)
	ConfigScoped(scope config.Scope) (*config.Config, error)
	CreateBranch(c *config.Branch) error
	CreateRemote(c *config.RemoteConfig) (*git.Remote, error)
	CreateRemoteAnonymous(c *config.RemoteConfig) (*git.Remote, error)
	CreateTag(name string, hash plumbing.Hash, opts *git.CreateTagOptions) (*plumbing.Reference, error)
	DeleteBranch(name string) error
	DeleteObject(hash plumbing.Hash) error
	DeleteRemote(name string) error
	DeleteTag(name string) error
	Fetch(o *git.FetchOptions) error
	FetchContext(ctx context.Context, o *git.FetchOptions) error
	Head() (*plumbing.Reference, error)
	Log(o *git.LogOptions) (object.CommitIter, error)
	Notes() (storer.ReferenceIter, error)
	Object(t plumbing.ObjectType, h plumbing.Hash) (object.Object, error)
	Objects() (*object.ObjectIter, error)
	Prune(opt git.PruneOptions) error
	Push(o *git.PushOptions) error
	PushContext(ctx context.Context, o *git.PushOptions) error
	Reference(name plumbing.ReferenceName, resolved bool) (*plumbing.Reference, error)
	References() (storer.ReferenceIter, error)
	Remote(name string) (*git.Remote, error)
	Remotes() ([]*git.Remote, error)
	RepackObjects(cfg *git.RepackConfig) (err error)
	ResolveRevision(rev plumbing.Revision) (*plumbing.Hash, error)
	SetConfig(cfg *config.Config) error
	Tag(name string) (*plumbing.Reference, error)
	TagObject(h plumbing.Hash) (*object.Tag, error)
	TagObjects() (*object.TagIter, error)
	Tags() (storer.ReferenceIter, error)
	TreeObject(h plumbing.Hash) (*object.Tree, error)
	TreeObjects() (*object.TreeIter, error)
	Worktree() (*git.Worktree, error)
}

type repositoryOption = func(repo *GitRepository)

type factoryMethod = func() (GoGitRepository, storage.Storer, error)
