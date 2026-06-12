package sandbox

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	ImageName    = "ricebox-sandbox"
	DockerfileDir = "docker"
)

type Sandbox struct {
	ContainerID string
	RicePath    string
}

func BuildImage() error {
	fmt.Println("🏗️  Building sandbox image...")
	cmd := exec.Command("docker", "build",
		"-t", ImageName,
		"-f", filepath.Join(DockerfileDir, "Dockerfile.sandbox"),
		DockerfileDir,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Start() (*Sandbox, error) {
	fmt.Println("🚀 Starting sandbox container...")

	cmd := exec.Command("docker", "run",
		"-d",
		"--rm",
		"-p", "5900:5900",
		"--name", "ricebox-sandbox",
		ImageName,
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	containerID := string(output)[:12]
	fmt.Printf("   Container ID: %s\n", containerID)
	fmt.Println("   VNC available at: localhost:5900")

	// Wait for container to be ready
	time.Sleep(2 * time.Second)

	return &Sandbox{ContainerID: containerID}, nil
}

func (s *Sandbox) Exec(command ...string) error {
	args := append([]string{"exec", "-it", "ricebox-sandbox"}, command...)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *Sandbox) CopyRiceToContainer(ricePath string) error {
	fmt.Printf("📤 Copying rice to container: %s\n", ricePath)
	cmd := exec.Command("docker", "cp",
		ricePath,
		"ricebox-sandbox:/home/rice/rices/",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Stop() error {
	fmt.Println("🛑 Stopping sandbox...")
	cmd := exec.Command("docker", "stop", "ricebox-sandbox")
	return cmd.Run()
}

func Cleanup() {
	exec.Command("docker", "rm", "-f", "ricebox-sandbox").Run()
}