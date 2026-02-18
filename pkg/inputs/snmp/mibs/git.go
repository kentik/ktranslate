package mibs

import (
	"context"
	"os"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
)

func cloneFromGit(ctx context.Context, profileDir string, gitUrl string, gitHash string, log logger.ContextL) error {
	var auth *githttp.BasicAuth
	if token := os.Getenv(snmp_util.KT_GIT_ACCESS_TOKEN); token != "" {
		auth = &githttp.BasicAuth{
			Username: os.Getenv(snmp_util.KT_GIT_ACCESS_USERNAME),
			Password: token,
		}
	}

	r, err := git.PlainCloneContext(ctx, profileDir, &git.CloneOptions{
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
