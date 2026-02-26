package registry

import (
	"fmt"
	"os"
	"os/exec"
)

// GitClone clones a git repository to destDir.
func GitClone(url, ref, destDir string) error {
	args := []string{"clone", "--branch", ref, "--single-branch", "--depth", "1", url, destDir}
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone %s: %w", url, err)
	}
	return nil
}

// GitPull runs git pull in the given repo directory.
func GitPull(repoDir string) error {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = repoDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git pull in %s: %w", repoDir, err)
	}
	return nil
}
