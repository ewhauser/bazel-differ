package internal

import (
	"os/exec"
)

type GitClient interface {
	Checkout(hash string) error
}

type gitClient struct {
	dir string
}

func NewGitClient(dir string) GitClient {
	return &gitClient{
		dir: dir,
	}
}

func (g gitClient) Checkout(hash string) error {
	cmd := exec.Command("git", "-C", g.dir, "checkout", hash, "--quiet")
	_, err := cmd.Output()
	return err
}
