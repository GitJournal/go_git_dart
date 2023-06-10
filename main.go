package main

import (
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func ExamplePlainClone() {

	url, directory, privateKeyFile := os.Args[1], os.Args[2], os.Args[3]
	var password string
	if len(os.Args) == 5 {
		password = os.Args[4]
	}

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

func main() {
	fmt.Println("Hello, playground")
	ExamplePlainClone()
}

// func ExamplePlainClone_usernamePassword() {
// 	// Tempdir to clone the repository
// 	dir, err := ioutil.TempDir("", "clone-example")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer os.RemoveAll(dir) // clean up

// 	// Clones the repository into the given dir, just as a normal git clone does
// 	_, err = git.PlainClone(dir, false, &git.CloneOptions{
// 		URL: "https://github.com/git-fixtures/basic.git",
// 		Auth: &http.BasicAuth{
// 			Username: "username",
// 			Password: "password",
// 		},
// 	})

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
