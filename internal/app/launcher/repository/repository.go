package repository

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Clone(repositoryURL string, cloneDepth int, commitHash string, localDir string) error {
	// Clone the repository
	repository, err := git.PlainClone(localDir, false, &git.CloneOptions{
		URL:   repositoryURL,
		Depth: cloneDepth,
	})
	if err != nil {
		return err
	}

	// Get the worktree
	workTree, err := repository.Worktree()
	if err != nil {
		return err
	}

	// Checkout the commit
	err = workTree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commitHash),
	})
	if err != nil {

		return err
	}

	return nil
}
