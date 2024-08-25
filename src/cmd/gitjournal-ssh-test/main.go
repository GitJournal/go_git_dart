package main

import (
	"fmt"
	"os"

	git "github.com/gitjournal/go-git-dart/internal/git"
	keygen "github.com/gitjournal/go-git-dart/internal/keygen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command: clone, fetch, push, defaultBranch")
		return
	}

	command := os.Args[1]
	switch command {
	case "clone":
		if len(os.Args) != 6 {
			fmt.Println("Usage: clone <url> <directory> <pemFile> <pemPassword>")
			return
		}
		pemBytes, err := os.ReadFile(os.Args[4])
		if err != nil {
			fmt.Println("Error reading PEM file:", err)
			return
		}
		err = git.Clone(os.Args[2], os.Args[3], pemBytes, os.Args[5])
		if err != nil {
			fmt.Println("Error cloning:", err)
			return
		}

	case "fetch", "push", "defaultBranch":
		if len(os.Args) != 5 {
			fmt.Printf("Usage: %s <remote> <pemFile> <pemPassword>\n", command)
			return
		}
		directory, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}
		fmt.Println(directory)
		pemBytes, err := os.ReadFile(os.Args[3])
		if err != nil {
			fmt.Println("Error reading PEM file:", err)
			return
		}

		switch command {
		case "fetch":
			err := git.Fetch(os.Args[2], directory, pemBytes, os.Args[4])
			if err != nil {
				fmt.Println("Error fetching:", err)
				return
			}

		case "push":
			err := git.Push(os.Args[2], directory, pemBytes, os.Args[4])
			if err != nil {
				fmt.Println("Error pushing:", err)
				return
			}

		case "defaultBranch":
			branch, err := git.DefaultBranch(os.Args[2], pemBytes, os.Args[4])
			if err != nil {
				fmt.Println("Error getting default branch:", err)
				return
			}
			fmt.Printf("DefaultBranch: %s\n", branch)
		}
	case "keygen":
		publicKey, privateKey, err := keygen.GenerateRSAKeys()
		if err != nil {
			fmt.Println("Error generating keys:", err)
			return
		}

		fmt.Printf("Public Key: %s\n", publicKey)
		fmt.Printf("Private Key: %s\n", privateKey)
		return

	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
