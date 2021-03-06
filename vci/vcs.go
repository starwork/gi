// Copyright (c) 2018, The GoGi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vci

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Masterminds/vcs"
	"github.com/goki/ki/dirs"
)

var (
	// ErrUnknownVCS is returned when VCS cannot be determined from the vcs Repo
	ErrUnknownVCS = errors.New("Unknown VCS")
)

// Repo provides an interface extending vcs.Repo
// (https://github.com/Masterminds/vcs)
// with support for file status information and operations.
type Repo interface {
	// vcs.Repo includes those interface functions
	vcs.Repo

	// Files returns a map of the current files and their status.
	Files() (Files, error)

	// Status returns status of given file -- returns Untracked and error
	// message on any error. FileStatus is a summary status category,
	// and string return value is more detailed status information formatted
	// according to standard conventions of given VCS.
	Status(fname string) (FileStatus, string)

	// Add adds the file to the repo
	Add(fname string) error

	// Move moves the file using VCS command to keep it updated
	Move(oldpath, newpath string) error

	// Delete removes the file from the repo and working copy.
	// Uses "force" option to ensure deletion.
	Delete(fname string) error

	// DeleteRemote removes the file from the repo but keeps the local file itself
	DeleteRemote(fname string) error

	// CommitFile commits a single file
	CommitFile(fname string, message string) error

	// RevertFile reverts a single file to the version that it was last in VCS,
	// losing any local changes (destructive!)
	RevertFile(fname string) error
}

func NewRepo(remote, local string) (Repo, error) {
	repo, err := vcs.NewRepo(remote, local)
	if err == nil {
		switch repo.Vcs() {
		case vcs.Git:
			r := &GitRepo{}
			r.GitRepo = *(repo.(*vcs.GitRepo))
			return r, err
		case vcs.Svn:
			r := &SvnRepo{}
			r.SvnRepo = *(repo.(*vcs.SvnRepo))
			return r, err
		case vcs.Hg:
			err = fmt.Errorf("Hg version control not yet supported")
		case vcs.Bzr:
			err = fmt.Errorf("Bzr version control not yet supported")
		}
	}
	return nil, err
}

// DetectRepo attemps to detect the presence of a repository at the given
// directory path -- returns type of repository if found, else vcs.NoVCS.
// Very quickly just looks for signature file name:
// .git for git
// .svn for svn -- but note that this will find any subdir in svn repo
func DetectRepo(path string) vcs.Type {
	if dirs.HasFile(path, ".git") {
		return vcs.Git
	}
	if dirs.HasFile(path, ".svn") {
		return vcs.Svn
	}
	// todo: rest later..
	return vcs.NoVCS
}

// RelPath return the path relative to the repository LocalPath()
func RelPath(repo Repo, path string) string {
	relpath, _ := filepath.Rel(repo.LocalPath(), path)
	return relpath
}
