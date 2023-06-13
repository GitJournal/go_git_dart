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

	_, err = git.PlainClone(directory, false, &git.CloneOptions{
		Auth:     publicKeys,
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("git clone failed:", err.Error())
		return errGitCloneFailed
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
