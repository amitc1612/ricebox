package rice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Applier struct {
	RicePath string
	Manifest *Manifest
	DryRun   bool
}

func NewApplier(ricePath string, dryRun bool) (*Applier, error) {
	m, err := LoadManifest(ricePath)
	if err != nil {
		return nil, err
	}

	return &Applier{
		RicePath: ricePath,
		Manifest: m,
		DryRun:   dryRun,
	}, nil
}

func (a *Applier) Apply() error {
	fmt.Printf("🍚 Applying rice: %s v%s\n", a.Manifest.Name, a.Manifest.Version)
	fmt.Printf("   %s\n\n", a.Manifest.Description)

	if err := a.applyDotfiles(); err != nil {
		return fmt.Errorf("dotfiles: %w", err)
	}

	if err := a.installPackages(); err != nil {
		return fmt.Errorf("packages: %w", err)
	}

	if err := a.restartServices(); err != nil {
		return fmt.Errorf("services: %w", err)
	}

	fmt.Println("\n✅ Rice applied successfully!")
	return nil
}

func (a *Applier) applyDotfiles() error {
	fmt.Println("🔗 Linking dotfiles...")
	for srcRel, destRel := range a.Manifest.Dotfiles {
		src := filepath.Join(a.RicePath, "dotfiles", srcRel)
		dest := expandPath(destRel)

		if a.DryRun {
			fmt.Printf("   [DRY RUN] %s → %s\n", src, dest)
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(dest), err)
		}

		// Remove existing
		if _, err := os.Lstat(dest); err == nil {
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("remove %s: %w", dest, err)
			}
		}

		// Create symlink
		if err := os.Symlink(src, dest); err != nil {
			return fmt.Errorf("symlink %s → %s: %w", src, dest, err)
		}
		fmt.Printf("   %s → %s\n", src, dest)
	}
	return nil
}

func (a *Applier) installPackages() error {
	if len(a.Manifest.Dependencies.Pacman) == 0 {
		return nil
	}

	fmt.Printf("📦 Installing packages: %s\n", strings.Join(a.Manifest.Dependencies.Pacman, ", "))
	if a.DryRun {
		return nil
	}

	args := append([]string{"pacman", "-S", "--needed", "--noconfirm"}, a.Manifest.Dependencies.Pacman...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a *Applier) restartServices() error {
	for _, service := range a.Manifest.Services.Restart {
		fmt.Printf("🔄 Restarting: %s\n", service)
		if a.DryRun {
			continue
		}
		cmd := exec.Command("systemctl", "--user", "restart", service)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("restart %s: %w", service, err)
		}
	}
	return nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}