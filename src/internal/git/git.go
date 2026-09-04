package git

import (
	"encoding/hex"
	"fmt"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"

	stdssh "golang.org/x/crypto/ssh"
)

func buildAuth(url string, privateKey []byte, password string) (transport.AuthMethod, error) {
	ep, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, err
	}

	publicKeys, err := ssh.NewPublicKeys(ep.User, privateKey, password)
	if err != nil {
		return nil, err
	}
	publicKeys.HostKeyCallback = stdssh.InsecureIgnoreHostKey()

	return publicKeys, nil
}

func Clone(url string, directory string, privateKey []byte, password string) error {
	auth, err := buildAuth(url, privateKey, password)
	if err != nil {
		return err
	}

	/*
		progressFile, err := os.OpenFile("/tmp/123.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		defer progressFile.Close()
	*/

	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		Auth: auth,
		URL:  url,
		// Progress: progressFile,
	})
	if err != nil {
		return err
	}

	return nil
}

func buildAuthForRemote(repo *git.Repository, remoteName string, privateKey []byte, password string) (transport.AuthMethod, error) {
	rem, err := repo.Remote(remoteName)
	if err != nil {
		return nil, err
	}

	urls := rem.Config().URLs
	if len(urls) == 0 {
		return nil, fmt.Errorf("no remote url")
	}

	return buildAuth(urls[0], privateKey, password)
}

func Fetch(remote string, directory string, privateKey []byte, password string) error {
	r, err := git.PlainOpen(directory)
	if err != nil {
		return err
	}

	auth, err := buildAuthForRemote(r, remote, privateKey, password)
	if err != nil {
		return err
	}

	err = r.Fetch(&git.FetchOptions{RemoteName: remote, Auth: auth})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func Pull(remote string, directory string, privateKey []byte, password string) error {
	return pull(Fetch, MergeCurrentBranch, remote, directory, privateKey, password)
}

func pull(
	fetchFn func(string, string, []byte, string) error,
	mergeFn func(string) error,
	remote string,
	directory string,
	privateKey []byte,
	password string,
) error {
	if err := fetchFn(remote, directory, privateKey, password); err != nil {
		return err
	}

	return mergeFn(directory)
}

func Push(remote string, directory string, privateKey []byte, password string) error {
	r, err := git.PlainOpen(directory)
	if err != nil {
		return err
	}

	auth, err := buildAuthForRemote(r, remote, privateKey, password)
	if err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{RemoteName: remote, Auth: auth})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func DefaultBranch(remoteUrl string, privateKey []byte, password string) (string, error) {
	auth, err := buildAuth(remoteUrl, privateKey, password)
	if err != nil {
		return "", err
	}

	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteUrl},
	})

	refs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		return "", err
	}

	defaultBranch := ""
	for _, ref := range refs {
		if ref.Name() == "HEAD" {
			defaultBranch = ref.Target().Short()
			break
		}
	}

	return defaultBranch, nil
}

func Add(directory string, path string) error {
	w, err := openWorktree(directory)
	if err != nil {
		return err
	}

	_, err = w.Add(path)
	return err
}

func Commit(directory string, message string) (string, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return "", fmt.Errorf("commit message is required")
	}

	r, w, err := openRepositoryAndWorktree(directory)
	if err != nil {
		return "", err
	}

	author, committer, err := commitSignatures(r)
	if err != nil {
		return "", err
	}

	hash, err := w.Commit(message, &git.CommitOptions{
		Author:    author,
		Committer: committer,
	})
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

func Remove(directory string, path string) error {
	w, err := openWorktree(directory)
	if err != nil {
		return err
	}

	_, err = w.Remove(path)
	return err
}

func Move(directory string, fromPath string, toPath string) error {
	w, err := openWorktree(directory)
	if err != nil {
		return err
	}

	_, err = w.Move(fromPath, toPath)
	return err
}

func ResetHard(directory string) error {
	_, w, err := openRepositoryAndWorktree(directory)
	if err != nil {
		return err
	}

	return w.Reset(&git.ResetOptions{Mode: git.HardReset})
}

func ResetHardTo(directory string, commitHash string) error {
	_, w, err := openRepositoryAndWorktree(directory)
	if err != nil {
		return err
	}

	commitHash = strings.TrimSpace(commitHash)
	if len(commitHash) != 40 {
		return fmt.Errorf("commit hash must be 40 hexadecimal characters")
	}
	if _, err := hex.DecodeString(commitHash); err != nil {
		return fmt.Errorf("invalid commit hash: %w", err)
	}

	return w.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(commitHash),
		Mode:   git.HardReset,
	})
}

func Switch(directory string, branch string) error {
	branch = strings.TrimSpace(branch)
	if branch == "" {
		return fmt.Errorf("branch name is required")
	}

	r, w, err := openRepositoryAndWorktree(directory)
	if err != nil {
		return err
	}

	branchRef := plumbing.NewBranchReferenceName(branch)
	if _, err := r.Reference(branchRef, true); err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{Branch: branchRef})
}

func openWorktree(directory string) (*git.Worktree, error) {
	_, w, err := openRepositoryAndWorktree(directory)
	return w, err
}

func openRepositoryAndWorktree(directory string) (*git.Repository, *git.Worktree, error) {
	r, err := git.PlainOpen(directory)
	if err != nil {
		return nil, nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, nil, err
	}

	return r, w, nil
}
