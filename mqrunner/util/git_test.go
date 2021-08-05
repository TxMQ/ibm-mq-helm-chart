package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCloneGitRepo(t *testing.T) {

	dir, err := ioutil.TempDir("", "git-clone")
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer func() { _ = os.RemoveAll(dir) }()

	url := "https://github.com/szesto/mq-operator.git"

	err = CloneGitRepo(dir, url)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("cloned git repo %s\n", url)
}
