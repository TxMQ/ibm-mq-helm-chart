package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCloneGitRepo(t *testing.T) {

	dir, err := ioutil.TempDir("", "git-fetch")
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	defer func() { _ = os.RemoveAll(dir) }()

	gitconfig := FetchConfig{
		Url:           "https://github.com/szesto/wiki.git",
		ReferenceName: "hello",
		Dir:           "",
	}

	err = GitFetch(dir, gitconfig)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fmt.Printf("fetched git repo %s to %s\n", gitconfig.Url, dir)
}

func TestFetchMergeConfigFiles(t *testing.T) {
	gitconfig := FetchConfig{
		Url:           "https://github.com/szesto/wiki.git",
		ReferenceName: "main",
		Dir:           "zosconn",
	}

	err := FetchMergeConfigFiles(gitconfig, GetMqscic(), GetQmini())
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
