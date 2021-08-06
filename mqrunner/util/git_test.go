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

	gitconfig := GitCloneConfig{
		Url:           "https://github.com/szesto/mq-operator.git",
		ReferenceName: "",
		Tag:           "",
		Dir:           "",
	}

	err = CloneGitRepo(dir, gitconfig)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("cloned git repo %s to %s\n", gitconfig.Url, dir)
}

func TestFetchMergeConfigFiles(t *testing.T) {
	gitconfig := GitCloneConfig{
		Url:           "https://github.com/szesto/mq-operator.git",
		ReferenceName: "",
		Tag:           "",
		Dir:           "values",
	}

	err := FetchMergeConfigFiles(gitconfig, GetMqscic(), GetQmini())
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
