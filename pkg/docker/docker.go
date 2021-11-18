package docker

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func Login(username, password string, repository string) error {
	loginRepo := strings.Split(repository, "/")[0]
	cmd := exec.Command("docker", "login", "-u", username, "--password-stdin", loginRepo)
	cmd.Stderr = os.Stderr

	stdin, _ := cmd.StdinPipe()
	cmd.Start()

	io.WriteString(stdin, password+"\n")
	return cmd.Wait()
}
