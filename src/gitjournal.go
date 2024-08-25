package main

import (
	git "github.com/gitjournal/go-git-dart/internal/git"
	keygen "github.com/gitjournal/go-git-dart/internal/keygen"

	/*
	   #include <stdlib.h>
	*/
	"C"
	"unsafe"
)

//export GitClone
func GitClone(url *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) *C.char {
	err := git.Clone(C.GoString(url), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != nil {
		return C.CString(err.Error())
	}

	return nil
}

//export GitFetch
func GitFetch(remote *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) *C.char {
	err := git.Fetch(C.GoString(remote), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != nil {
		return C.CString(err.Error())
	}

	return nil
}

//export GitPush
func GitPush(remote *C.char, directory *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char) *C.char {
	err := git.Push(C.GoString(remote), C.GoString(directory), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != nil {
		return C.CString(err.Error())
	}

	return nil
}

//export GitDefaultBranch
func GitDefaultBranch(remoteUrl *C.char, privateKey *C.char, privateKeyLen C.int, password *C.char, outputBranchName **C.char) *C.char {
	val, err := git.DefaultBranch(C.GoString(remoteUrl), C.GoBytes(unsafe.Pointer(privateKey), privateKeyLen), C.GoString(password))
	if err != nil {
		return C.CString(err.Error())
	}

	*outputBranchName = C.CString(val)
	return nil
}

//export GJGenerateRSAKeys
func GJGenerateRSAKeys(publicKey **C.char, privateKey **C.char) *C.char {
	publicKeyVal, privateKeyVal, err := keygen.GenerateRSAKeys()
	if err != nil {
		return C.CString(err.Error())
	}

	*publicKey = C.CString(publicKeyVal)
	*privateKey = C.CString(privateKeyVal)

	return nil
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
