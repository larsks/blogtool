package git

import (
	"os/exec"
	"strings"
)

func GetTopdir() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(out), "\n"), err
}

func CreateBranch(branch, start string) error {
	if start == "" {
		start = "master"
	}

	if err := exec.Command("git", "checkout", "-b", branch, start).Run(); err != nil {
		return err
	}

	return nil
}

func AddFiles(files ...string) error {
	args := []string{"add"}
	args = append(args, files...)

	if err := exec.Command("git", args...).Run(); err != nil {
		return err
	}

	return nil
}

func Commit(message string) error {
	if err := exec.Command("git", "commit", "-m", message).Run(); err != nil {
		return err
	}

	return nil
}
