package main

import (

	/*
	   #include <stdlib.h>
	*/
	"C"
	"fmt"
	"unsafe"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"

	stdssh "golang.org/x/crypto/ssh"
)

//export GitClone
func GitClone(url *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) *C.char {
	err := gitClone(C.GoString(url), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != nil {
		return C.CString(err.Error())
	}

	return nil
}

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

func gitClone(url string, directory string, privateKey []byte, password string) error {
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

//export GitFetch
func GitFetch(remote *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) *C.char {
	err := gitFetch(C.GoString(remote), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != nil {
		return C.CString(err.Error())
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

func gitFetch(remote string, directory string, privateKey []byte, password string) error {
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

//export GitPush
func GitPush(remote *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) *C.char {
	err := gitPush(C.GoString(remote), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != nil {
		return C.CString(err.Error())
	}

	return nil
}

func gitPush(remote string, directory string, privateKey []byte, password string) error {
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

/*
type GitDefaultBranchResult struct {
	err int
	val *C.char
}
*/

//export GitDefaultBranch
func GitDefaultBranch(remoteUrl *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) *C.char {
	err, val := gitDefaultBranch(C.GoString(remoteUrl), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != 0 {
		return nil
	}
	return C.CString(val)
}

func gitDefaultBranch(remoteUrl string, privateKey []byte, password string) (int, string) {
	auth, err := buildAuth(remoteUrl, privateKey, password)
	if err != nil {
		return 1, ""
	}

	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteUrl},
	})

	refs, err := remote.List(&git.ListOptions{Auth: auth})
	if err != nil {
		fmt.Println("git remote list failed:", err.Error())
		return 1, ""
	}

	defaultBranch := ""
	for _, ref := range refs {
		if ref.Name() == "HEAD" {
			defaultBranch = ref.Target().Short()
			break
		}
	}

	return 0, defaultBranch
}

func main() {}

/*
func main() {
	fmt.Println("Hello, playground")

	url, directory, privateKeyFile := os.Args[1], os.Args[2], os.Args[3]
	var password string
	if len(os.Args) == 5 {
		password = os.Args[4]
	}

	privateKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("URL:", url)
	fmt.Println("Directory:", directory)
	fmt.Println("PrivateKey:", privateKey)
	fmt.Println("Password:", password)

	err = gitClone(url, directory, privateKey, password)
	if err != nil {
		panic(err)
	}
}
*/
