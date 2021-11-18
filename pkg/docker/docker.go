package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Login(username, password string, repository string) error {
	loginRegistry := strings.Split(repository, "/")[0]
	fmt.Fprintf(os.Stderr, "login registry: %s\n", loginRegistry)
	cmd := exec.Command("docker", "login", "-u", username, "-p", password, loginRegistry)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
