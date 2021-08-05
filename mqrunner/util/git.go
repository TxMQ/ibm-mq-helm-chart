package util

import (
	git "github.com/go-git/go-git/v5"
)

// Clone git repository. There could be more parameters
// eg auth, ref, etc

func CloneGitRepo(path string, url string) error {

	options := git.CloneOptions {
		URL:               url, // https://github.com/.../foo.git
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
