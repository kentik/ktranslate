package mibs

import (
	"context"
	"io"
	"os"

	"github.com/kentik/ktranslate/pkg/eggs/logger"
	snmp_util "github.com/kentik/ktranslate/pkg/inputs/snmp/util"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

func cloneFromGit(ctx context.Context, profileDir string, gitUrl string, gitHash string, log logger.ContextL) error {
	var branch plumbing.ReferenceName
	if bb := os.Getenv(snmp_util.KT_GIT_PULL_BRANCH); bb != "" {
		branch = plumbing.NewBranchReferenceName(bb)
	}

	cloneOpts := &git.CloneOptions{
		URL:      gitUrl,
		Auth:     snmp_util.GetGitCreds(),
		Progress: io.Discard,
	}
	// If a branch is specified, clone that branch directly instead of
	// cloning the default branch and manually rewriting references.
	if branch != "" {
		cloneOpts.ReferenceName = branch
		cloneOpts.SingleBranch = true
		log.Infof("Checking profile repo out to branch %s", branch)
	}
	r, err := git.PlainCloneContext(ctx, profileDir, cloneOpts)
	if err != nil {
		return err
	}

	if gitHash == "" {
		return nil
	}

	log.Infof("Checking profile repo out to hash %s ", gitHash)
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	coOpts := &git.CheckoutOptions{
		Hash: plumbing.NewHash(gitHash),
	}
	err = w.Checkout(coOpts)
	return err
}
