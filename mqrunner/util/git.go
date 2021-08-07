package util

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

type FetchConfig struct {
	Url string	// git repo url
	ReferenceName string
	Dir string // directory in repository
}

func GitFetch(path string, config FetchConfig) error {

	// git init path
	if out, err := git("init", path); err != nil {
		return err
	} else {
		log.Printf("git-fetch: %s\n", out)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	if err = os.Chdir(path); err != nil {
		return err
	}

	defer func() {_ = os.Chdir(cwd)}()

	// git remote add origin url
	if out, err := git("remote", "add", "origin", config.Url); err != nil {
		return err
	} else {
		log.Printf("git-fetch: %s\n", out)
	}

	// git config http.sslVerify false
	if out, err := git("config", "http.sslVerify", "false"); err != nil {
		return err
	} else {
		log.Printf("git-fetch: %s\n", out)
	}

	rev := ""
	if len(config.ReferenceName) == 0 {
		rev = "HEAD"
		// git symbolic-ref HEAD refs/remotes/origin/HEAD
		if out, err := git("symbolic-ref", "HEAD", "refs/remotes/origin/HEAD"); err != nil {
			return err
		} else {
			log.Printf("git-fetch: %s\n", out)
		}
	} else {
		rev = config.ReferenceName
	}

	// git fetch --update-head-ok --force referenceName
	if out, err := git("fetch", "origin", "--update-head-ok", "--force", rev); err != nil {
		return err
	} else {
		log.Printf("git-fetch: %s\n", out)
	}

	// git checkout fetch_head
	if out, err := git("checkout", "FETCH_HEAD"); err != nil {
		return err
	} else {
		log.Printf("git-fetch: %s\n", out)
	}

	return nil
}

func git(args ...string) (string, error) {
	c := exec.Command("/usr/bin/git", args...)

	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out

	if err := c.Run(); err != nil {
		return "", err
	}

	return c.String(), nil
}