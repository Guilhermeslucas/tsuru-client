// Copyright 2012 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"github.com/fsouza/gogit/git"
	"regexp"
)

// AppGuesser is used to guess the name of an app based in a file path.
type AppGuesser interface {
	GuessName(path string) (string, error)
}

// GitGuesser uses git to guess the name of the app.
//
// It reads the "tsuru" remote from git config file. If the remote does not
// exist, or does not match the tsuru pattern (git@<something>:<app-name>.git),
// GuessName will return an error.
type GitGuesser struct{}

func (g GitGuesser) GuessName(path string) (string, error) {
	repoPath, err := git.DiscoverRepositoryPath(path)
	if err != nil {
		return "", fmt.Errorf("Git repository not found: %s.", err)
	}
	repo, err := git.OpenRepository(repoPath)
	if err != nil {
		return "", err
	}
	defer repo.Free()
	config, err := repo.Config()
	if err != nil {
		return "", err
	}
	remoteUrl, err := config.GetString("remote.tsuru.url")
	if err != nil {
		return "", errors.New("tsuru remote not declared.")
	}
	re := regexp.MustCompile(`^git@.*:(.*)\.git$`)
	matches := re.FindStringSubmatch(remoteUrl)
	if len(matches) < 2 {
		return "", fmt.Errorf(`"tsuru" remote did not match the pattern. Want something like git@<host>:<app-name>.git, got %s`, remoteUrl)
	}
	return matches[1], nil
}
