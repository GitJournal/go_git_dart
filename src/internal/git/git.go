package git

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
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
