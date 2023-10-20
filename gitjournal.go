package main

import (
	"os"

	"C"
	"unsafe"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)
import "fmt"

const errPublicKeysFailed = 55
const errGitCloneFailed = 56
const errGitOpenFailed = 57
const errGitPullFailed = 59

//export GitClone
func GitClone(url *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) int {
	return gitClone(C.GoString(url), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
}

func gitClone(url string, directory string, privateKey []byte, password string) int {
	publicKeys, err := ssh.NewPublicKeys("git", privateKey, password)
	if err != nil {
		fmt.Println("generate publickeys failed:", err.Error())
		return errPublicKeysFailed
	}

	progressFile, err := os.OpenFile("/tmp/123.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer progressFile.Close()

	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		Auth:     publicKeys,
		URL:      url,
		Progress: progressFile,
	})
	if err != nil {
		fmt.Println("git clone failed:", err.Error())
		return errGitCloneFailed
	}

	return 0
}

//export GitFetch
func GitFetch(remote *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) int {
	return gitFetch(C.GoString(remote), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
}

func gitFetch(remote string, directory string, privateKey []byte, password string) int {
	publicKeys, err := ssh.NewPublicKeys("git", privateKey, password)
	if err != nil {
		fmt.Println("generate publickeys failed:", err.Error())
		return errPublicKeysFailed
	}

	fmt.Println("git fetch", directory)
	r, err := git.PlainOpen(directory)
	if err != nil {
		fmt.Println("git open failed:", err.Error())
		return errGitOpenFailed
	}

	err = r.Fetch(&git.FetchOptions{RemoteName: remote, Auth: publicKeys})
	if err == git.NoErrAlreadyUpToDate {
		return 0
	}

	if err != nil {
		fmt.Println("git pull failed:", err.Error())
		return errGitPullFailed
	}

	return 0
}

//export GitPush
func GitPush(remote *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) int {
	return gitPush(C.GoString(remote), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
}

func gitPush(remote string, directory string, privateKey []byte, password string) int {
	publicKeys, err := ssh.NewPublicKeys("git", privateKey, password)
	if err != nil {
		fmt.Println("generate publickeys failed:", err.Error())
		return errPublicKeysFailed
	}

	fmt.Println("git push", directory)
	r, err := git.PlainOpen(directory)
	if err != nil {
		fmt.Println("git open failed:", err.Error())
		return errGitOpenFailed
	}

	err = r.Push(&git.PushOptions{RemoteName: remote, Auth: publicKeys})
	if err == git.NoErrAlreadyUpToDate {
		return 0
	}

	if err != nil {
		fmt.Println("git push failed:", err.Error())
		return errGitPullFailed
	}

	return 0
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
