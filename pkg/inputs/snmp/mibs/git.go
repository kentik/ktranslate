package mibs

import (
	"context"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func cloneFromGit(ctx context.Context, profileDir string, gitUrl string, gitHash string, log logger.ContextL) error {
	r, err := git.PlainCloneContext(ctx, profileDir, false, &git.CloneOptions{
		URL: gitUrl,
	})
	if err != nil {
		return err
	}
	if gitHash == "" {
		return nil
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	log.Infof("Checking profile repo out to %s", gitHash)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(gitHash),
	})
	return err
}
