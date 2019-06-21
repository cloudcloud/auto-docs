// Package data provides functionality to interact with sources
// of documentation data.
//
// In the initial phase, these functions are quite... poor. Some
// degree of effort needs investing to better process the steps
// necessary to work with Git. There are several steps required
// that each generate their own errors, care is required.
package data

import (
	"log"
	"time"

	autodocs "github.com/cloudcloud/auto-docs"
	"github.com/cloudcloud/auto-docs/auto-docs/docs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// State will keep a hold of the current details surrounding the
// local data repository.
type State struct {
	// URL denotes the remote for this State.
	Git autodocs.Git

	// Sha contains the current sha checked out.
	Sha string

	// g holds the git handler instance
	g *git.Repository
}

// Prep will setup and begin the data acquisition for a specific
// repository location.
func Prep(g autodocs.Git) *State {
	s := &State{
		Git: g,
	}
	auth, err := ssh.NewPublicKeysFromFile(g.Username, g.SSHKey, g.Password)
	if err != nil {
		log.Fatalf("unable to setup ssh: %s\n", err)
	}

	_, err = git.PlainClone(
		g.LocalPath,
		false,
		&git.CloneOptions{
			Auth:          auth,
			URL:           g.URI,
			ReferenceName: plumbing.NewBranchReferenceName(g.Branch),
			Depth:         1,
			SingleBranch:  true,
		},
	)
	if err != nil && err.Error() != "repository already exists" {
		log.Println("unable to clone repository:", err)
	}

	s.g, err = git.PlainOpen(g.LocalPath)
	if err != nil {
		log.Println("unable to open repository:", err)
	}

	sha, err := s.g.Head()
	if err != nil {
		log.Println("unable to retrieve commit:", err)
	}
	s.Sha = sha.String()

	// tell data to update
	docs.S.UpdateFromPath(g.LocalPath)

	return s
}

// Fetch will check for a new sha and pull if one found.
func (s *State) Fetch(t time.Time) {
	// pull from the upstream
	w, err := s.g.Worktree()
	if err != nil {
		log.Println("unable to get worktree:", err)
	}
	auth, err := ssh.NewPublicKeysFromFile(s.Git.Username, s.Git.SSHKey, s.Git.Password)
	if err != nil {
		log.Fatalf("unable to setup ssh: %s\n", err)
	}

	err = w.Pull(
		&git.PullOptions{
			Auth:          auth,
			Depth:         1,
			ReferenceName: plumbing.NewBranchReferenceName(s.Git.Branch),
			SingleBranch:  true,
		},
	)
	if err != nil && err.Error() != "already up-to-date" {
		log.Println("couldn't fetch:", err)
	}

	// compare the sha against the previous known hash
	if sha, _ := s.g.Head(); sha.String() != s.Sha {
		// tell data to re-process
		log.Println("updating local repo at", t)

		s.Sha = sha.String()
		docs.S.UpdateFromPath(s.Git.LocalPath)
	}
}
