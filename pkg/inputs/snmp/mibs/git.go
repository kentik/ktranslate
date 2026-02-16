package mibs

import (
	"context"
	"os"

	"github.com/kentik/ktranslate/pkg/eggs/logger"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
)

func cloneFromGit(ctx context.Context, profileDir string, gitUrl string, gitHash string, log logger.ContextL) error {
	var auth *githttp.BasicAuth
	if token := os.Getenv("KT_GITHUB_ACCESS_TOKEN"); token != "" {
		auth = &githttp.BasicAuth{
			Username: "foo", // yes, this can be anything except an empty string
			Password: token,
		}
	}

	r, err := git.PlainCloneContext(ctx, profileDir, false, &git.CloneOptions{
		URL:  gitUrl,
		Auth: auth,
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
