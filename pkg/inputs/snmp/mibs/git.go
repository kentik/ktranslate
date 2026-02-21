package mibs

import (
	"context"
	"os"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

func cloneFromGit(ctx context.Context, profileDir string, gitUrl string, gitHash string, log logger.ContextL) error {
	r, err := git.PlainCloneContext(ctx, profileDir, &git.CloneOptions{
		URL:  gitUrl,
		Auth: snmp_util.GetGitCreds(),
	})
	if err != nil {
		return err
	}

	var branch plumbing.ReferenceName
	if bb := os.Getenv(snmp_util.KT_GIT_PULL_BRANCH); bb != "" {
		branch = plumbing.NewBranchReferenceName(bb)
	}

	if gitHash == "" && branch == "" {
		return nil
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	var coOpts *git.CheckoutOptions
	if branch != "" {
		headRef, err := r.Head()
		if err != nil {
			return err
		}
		ref := plumbing.NewHashReference(branch, headRef.Hash())
		err = r.Storer.SetReference(ref)
		if err != nil {
			return err
		}
		coOpts = &git.CheckoutOptions{
			Branch: ref.Name(),
		}
		log.Infof("Checking profile repo out to branch %s", branch)
	} else {
		coOpts = &git.CheckoutOptions{
			Hash: plumbing.NewHash(gitHash),
		}
		log.Infof("Checking profile repo out to hash %s ", gitHash)
	}

	err = w.Checkout(coOpts)
	return err
}
