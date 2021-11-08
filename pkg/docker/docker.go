package docker

import (
	"os"
	"os/exec"
)

func Login(username, password string, repository string) error {
	cmd := exec.Command("docker", "login", "-u", username, "-p", password, repository)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
