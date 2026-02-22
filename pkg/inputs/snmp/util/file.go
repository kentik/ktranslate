package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
)

const (
	KT_GIT_ACCESS_TOKEN    = "KT_GIT_ACCESS_TOKEN"
	KT_GIT_ACCESS_USERNAME = "KT_GIT_ACCESS_USERNAME"
	KT_GIT_PULL_BRANCH     = "KT_GIT_PULL_BRANCH"
	KT_GIT_PUSH_BRANCH     = "KT_GIT_PUSH_BRANCH"
	KT_GIT_COMMIT_EMAIL    = "KT_GIT_COMMIT_EMAIL"
	KT_GIT_COMMIT_NAME     = "KT_GIT_COMMIT_NAME"
)

// Utility to load a config file from various places
func LoadFile(ctx context.Context, file string) ([]byte, error) {
	u, err := url.Parse(file)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		return loadFromHttp(ctx, file)
	case "s3":
		return loadFromS3(ctx, u, getS3Downloader())
	case "s3m":
		return loadFromS3(ctx, u, getMockS3Client())
	case "git":
		return loadFromGit(ctx, u)
	default:
		return os.ReadFile(file)
	}
}

// Utility to write a config file to various places
func WriteFile(ctx context.Context, file string, payload []byte, perms fs.FileMode) error {
	u, err := url.Parse(file)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "http", "https":
		return writeToHttp(ctx, file, payload)
	case "s3":
		return writeToS3(ctx, u, payload, getS3Uploader())
	case "s3m":
		return writeToS3(ctx, u, payload, getMockS3Client())
	case "git":
		return writeToGit(ctx, u, payload, perms)
	default:
		return os.WriteFile(file, payload, perms)
	}
}

func loadFromHttp(ctx context.Context, file string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, file, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func writeToHttp(ctx context.Context, file string, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, file, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Cannot send to %s: %d", file, resp.StatusCode)
	}
	return nil
}

func loadFromS3(ctx context.Context, url *url.URL, client s3ClientDown) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	size, err := client.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(url.Host),
		Key:    aws.String(url.Path),
	})
	if err != nil {
		return nil, err
	}

	bufr := buf.Bytes()
	return bufr[0:size], nil
}

func writeToS3(ctx context.Context, url *url.URL, payload []byte, client s3ClientUp) error {
	_, err := client.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(url.Host),
		Body:   bytes.NewBuffer(payload),
		Key:    aws.String(url.Path),
	})
	return err
}

func writeToGit(ctx context.Context, url *url.URL, payload []byte, perms fs.FileMode) error {
	// Get repo cloned
	dir, err := os.MkdirTemp("", "ktrans")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up

	// If we are in a branch, use this here.
	var branch plumbing.ReferenceName
	var refSpecSet []config.RefSpec
	if bb := os.Getenv(KT_GIT_PUSH_BRANCH); bb != "" {
		branch = plumbing.NewBranchReferenceName(bb)
		refSpecSet = []config.RefSpec{config.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/heads/%s", bb, bb))}
	}

	// Clone things locally to kick off.
	filePath, r, err := gitClone(ctx, url, dir, branch)
	if err != nil {
		return err
	}

	file := path.Join(dir, filePath)

	// Copy new file onto path
	err = os.WriteFile(file, payload, perms)
	if err != nil {
		return fmt.Errorf("%s, WriteFile %s", err.Error(), file)
	}

	// Commit repo
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("%s, Worktree %s", err.Error(), file)
	}

	// Adds the new file to the staging area.
	_, err = w.Add(filePath)
	if err != nil {
		return fmt.Errorf("%s, git add %s", err.Error(), file)
	}

	name := os.Getenv(KT_GIT_COMMIT_NAME)
	if name == "" {
		name = "Ktranslate internal"
	}
	email := os.Getenv(KT_GIT_COMMIT_EMAIL)
	if email == "" {
		email = "ktranslate@kentik.com"
	}

	_, err = w.Commit("ktranslate adding new version of config file", &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("%s, git commit %s", err.Error(), file)
	}

	// Push repo.
	return r.Push(&git.PushOptions{
		Auth:     GetGitCreds(),
		Force:    true,
		RefSpecs: refSpecSet,
	})
}

func GetGitCreds() *githttp.BasicAuth {
	token := os.Getenv(KT_GIT_ACCESS_TOKEN)
	if token == "" {
		return nil
	}
	username := os.Getenv(KT_GIT_ACCESS_USERNAME)
	if username == "" {
		// Many Git hosting services require a non-empty username when using personal access tokens.
		// Default to a conventional placeholder username if none is provided via environment.
		username = "git"
	}
	return &githttp.BasicAuth{
		Username: username,
		Password: token,
	}
}

func gitClone(ctx context.Context, url *url.URL, dir string, branch plumbing.ReferenceName) (string, *git.Repository, error) {
	// Derive gitRepo from host and the first two path segments (owner/repo),
	// trimming an optional ".git" suffix from the repo name.
	cleanPath := strings.TrimPrefix(url.Path, "/")
	segments := strings.Split(cleanPath, "/")
	if len(segments) < 3 {
		return "", nil, fmt.Errorf("invalid git url path: %s. Expected format: git://host/owner/repo/path/to/file.yaml", url.String())
	}
	owner := segments[0]
	repo := segments[1]
	if strings.HasSuffix(repo, ".git") {
		repo = strings.TrimSuffix(repo, ".git")
	}
	gitRepo := "https://" + path.Join(url.Host, owner, repo) + ".git"
	filePath := filepath.Clean(path.Join(segments[2:]...))

	cloneOpts := &git.CloneOptions{
		URL:      gitRepo,
		Auth:     GetGitCreds(),
		Progress: io.Discard,
	}
	// If a branch is specified, clone that branch directly instead of
	// cloning the default branch and manually rewriting references.
	if branch != "" {
		cloneOpts.ReferenceName = branch
		cloneOpts.SingleBranch = true
	}
	r, err := git.PlainCloneContext(ctx, dir, cloneOpts)
	if err != nil {
		return "", nil, fmt.Errorf("failed to clone git repository %s: %w", gitRepo, err)
	}

	return filePath, r, err
}

func loadFromGit(ctx context.Context, url *url.URL) ([]byte, error) {
	dir, err := os.MkdirTemp("", "ktrans")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) // clean up

	var branch plumbing.ReferenceName
	if bb := os.Getenv(KT_GIT_PULL_BRANCH); bb != "" {
		branch = plumbing.NewBranchReferenceName(bb)
	}
	filePath, _, err := gitClone(ctx, url, dir, branch)
	if err != nil {
		return nil, err
	}

	file := path.Join(dir, filePath)
	return os.ReadFile(file)
}

type s3ClientDown interface {
	DownloadWithContext(aws.Context, io.WriterAt, *s3.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
}

type s3ClientUp interface {
	UploadWithContext(aws.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type s3Full interface {
	s3ClientUp
	s3ClientDown
}

func getS3Downloader() s3ClientDown {
	sess := session.Must(session.NewSession())
	client := s3manager.NewDownloader(sess)
	return client
}

func getS3Uploader() s3ClientUp {
	sess := session.Must(session.NewSession())
	client := s3manager.NewUploader(sess)
	return client
}

type mockS3 struct {
	lastContent []byte
}

var (
	sMock *mockS3
)

func (m *mockS3) UploadWithContext(ctx aws.Context, in *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	m.lastContent = in.Body.(*bytes.Buffer).Bytes()
	return nil, nil
}

func (m *mockS3) DownloadWithContext(ctx aws.Context, w io.WriterAt, in *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	w.WriteAt(m.lastContent, 0)
	return int64(len(m.lastContent)), nil
}

func getMockS3Client() s3Full {
	return sMock
}

func TouchFile(path string) error {
	u, err := url.Parse(path)
	if err != nil {
		return err
	}

	// If an external scheme, just return nil.
	switch u.Scheme {
	case "http", "https":
		return nil
	case "s3":
		return nil
	case "s3m":
		return nil
	case "git":
		return nil
	}

	_, err = os.Stat(path)
	if err != nil {
		return err
	}

	// Now see if can write.
	currentTime := time.Now().Local()
	return os.Chtimes(path, currentTime, currentTime)
}
