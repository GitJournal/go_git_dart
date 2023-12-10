package main

import (

	/*
	   #include <stdlib.h>
	*/
	"C"
	"fmt"
	"unsafe"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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

func gitClone(url string, directory string, privateKey []byte, password string) error {
	publicKeys, err := ssh.NewPublicKeys("git", privateKey, password)
	if err != nil {
		return err
	}
	publicKeys.HostKeyCallback = stdssh.InsecureIgnoreHostKey()

	/*
		progressFile, err := os.OpenFile("/tmp/123.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		defer progressFile.Close()
	*/

	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		Auth: publicKeys,
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

func gitFetch(remote string, directory string, privateKey []byte, password string) error {
	publicKeys, err := ssh.NewPublicKeys("git", privateKey, password)
	if err != nil {
		return err
	}
	publicKeys.HostKeyCallback = stdssh.InsecureIgnoreHostKey()

	fmt.Println("git fetch", directory)
	r, err := git.PlainOpen(directory)
	if err != nil {
		return err
	}

	err = r.Fetch(&git.FetchOptions{RemoteName: remote, Auth: publicKeys})
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
	publicKeys, err := ssh.NewPublicKeys("git", privateKey, password)
	if err != nil {
		return err
	}
	publicKeys.HostKeyCallback = stdssh.InsecureIgnoreHostKey()

	r, err := git.PlainOpen(directory)
	if err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{RemoteName: remote, Auth: publicKeys})
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
func GitDefaultBranch(remote *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) *C.char {
	err, val := gitDefaultBranch(C.GoString(remote), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != 0 {
		return nil
	}
	return C.CString(val)
}

func gitDefaultBranch(remoteName string, directory string, privateKey []byte, password string) (int, string) {
	publicKeys, err := ssh.NewPublicKeys("git", privateKey, password)
	if err != nil {
		fmt.Println("generate publickeys failed:", err.Error())
		return 1, ""
	}
	publicKeys.HostKeyCallback = stdssh.InsecureIgnoreHostKey()

	repo, err := git.PlainOpen(directory)
	if err != nil {
		fmt.Println("git open failed:", err.Error())
		return 1, ""
	}

	remote, err := repo.Remote(remoteName)
	if err != nil {
		fmt.Println("git remote failed:", err.Error())
		return 1, ""
	}

	refs, err := remote.List(&git.ListOptions{Auth: publicKeys})
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

	gitClone(url, directory, privateKeyFile, password)
}
*/
