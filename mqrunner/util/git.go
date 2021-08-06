package util

import (
	git "github.com/go-git/go-git/v5"
)

type GitCloneConfig struct {
	Url string	// git repo url
	ReferenceName string
	Tag string
	Dir string // directory in repository
}

// Clone git repository. There could be more parameters
// eg auth, ref, etc

func CloneGitRepo(path string, config GitCloneConfig) error {

	options := git.CloneOptions {
		URL:               config.Url, // https://github.com/.../foo.git
		Auth:              nil,
		RemoteName:        "",
		ReferenceName:     "",
		SingleBranch:      false,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: 0,
		Progress:          nil,
		Tags:              0,
		InsecureSkipTLS:   true,
		CABundle:          nil,
	}

	_, err := git.PlainClone(path, false, &options)
	if err != nil {
		return err
	}

	return nil
}
