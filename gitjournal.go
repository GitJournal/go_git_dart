package main

import (
	"fmt"
	"log"
	"os"

	"C"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

//export GitClone
func GitClone(url, directory, privateKeyFile, password *C.char) {
	gitCloneInternal(C.GoString(url), C.GoString(directory), C.GoString(privateKeyFile), C.GoString(password))
}

func gitCloneInternal(url, directory, privateKeyFile, password string) {
	_, err := os.Stat(privateKeyFile)
	if err != nil {
		log.Fatalln("read file", privateKeyFile, err.Error())
		return
	}

	// Clone the given repository to the given directory
	fmt.Println("git clone", url)
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		log.Fatalln("generate publickeys failed:", err.Error())
		return
	}

	fmt.Println("git clone", url, "to", directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		Auth:     publicKeys,
		URL:      url,
		Progress: os.Stdout,
	})
	fmt.Println(r, err)
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

	gitCloneInternal(url, directory, privateKeyFile, password)
}
*/
