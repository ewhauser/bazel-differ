package internal

import (
	"os/exec"
)

type GitClient interface {
	Checkout(hash string) string, error
}

type gitClient struct {
	dir string
}

func NewGitClient(dir string) GitClient {
	return &gitClient{
		dir: dir,
	}
}

func (g gitClient) Checkout(hash string) string, error {
	cmd := exec.Command("git", "-C", g.dir, "checkout", hash, "--quiet")
	output, err := cmd.CombinedOutput()
	return err
}
