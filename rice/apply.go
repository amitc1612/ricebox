package rice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Backup struct {
	Path      string
	Timestamp string
	Files     map[string]string
}

func Apply(ricePath string, dryRun bool) error {
	manifest, err := LoadManifest(ricePath)
	if err != nil {
		return err
	}

	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	fmt.Printf("\n🍚 Applying rice: %s v%s\n", manifest.Name, manifest.Version)
	if manifest.Description != "" {
		fmt.Printf("   %s\n", manifest.Description)
	}
	fmt.Println()

	if dryRun {
		fmt.Println("⚠ DRY RUN - no changes will be made\n")
	}

	backup, err := backupCurrent(manifest, dryRun)
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	if err := applyDotfiles(ricePath, manifest, dryRun); err != nil {
		fmt.Println("\n❌ Dotfiles failed! Restoring backup...")
		restoreBackup(backup)
		return fmt.Errorf("dotfiles: %w", err)
	}

	if err := installPackages(manifest, dryRun); err != nil {
		fmt.Println("\n⚠ Package installation had errors (dotfiles still applied)")
	}

	if err := restartServices(manifest, dryRun); err != nil {
		fmt.Println("\n⚠ Service restart had errors (dotfiles still applied)")
	}

	fmt.Printf("\n✅ Rice applied successfully!")
	if !dryRun {
		fmt.Printf(" Backup saved at: %s\n", backup.Path)
	}
	return nil
}

func backupCurrent(manifest *Manifest, dryRun bool) (*Backup, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(home, ".local/share/ricebox/backups", timestamp)

	b := &Backup{
		Path:      backupDir,
		Timestamp: timestamp,
		Files:     make(map[string]string),
	}

	for _, destRel := range manifest.Dotfiles {
		dest := expandPath(destRel)
		if _, err := os.Lstat(dest); os.IsNotExist(err) {
			continue
		}

		backupDest := filepath.Join(backupDir, destRel)
		b.Files[destRel] = dest

		if dryRun {
			fmt.Printf("   [DRY RUN] Would backup: %s\n", dest)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(backupDest), 0755); err != nil {
			return b, fmt.Errorf("mkdir for backup: %w", err)
		}

		if err := copyFile(dest, backupDest); err != nil {
			return b, fmt.Errorf("backup %s: %w", dest, err)
		}
		fmt.Printf("   📦 Backed up: %s\n", dest)
	}

	if !dryRun {
		fmt.Printf("   Backup location: %s\n", backupDir)
	}
	return b, nil
}

func applyDotfiles(ricePath string, manifest *Manifest, dryRun bool) error {
	fmt.Println("\n🔗 Linking dotfiles...")

	for srcRel, destRel := range manifest.Dotfiles {
		src := filepath.Join(ricePath, "dotfiles", srcRel)
		dest := expandPath(destRel)

		if _, err := os.Stat(src); os.IsNotExist(err) {
			return fmt.Errorf("source file not found: %s", src)
		}

		if dryRun {
			fmt.Printf("   [DRY RUN] Would link: %s → %s\n", src, dest)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(dest), err)
		}

		if _, err := os.Lstat(dest); err == nil {
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("remove existing %s: %w", dest, err)
			}
		}

		if err := os.Symlink(src, dest); err != nil {
			return fmt.Errorf("symlink %s → %s: %w", src, dest, err)
		}
		fmt.Printf("   🔗 %s → %s\n", src, dest)
	}
	return nil
}

func installPackages(manifest *Manifest, dryRun bool) error {
	if len(manifest.Dependencies.Pacman) == 0 {
		return nil
	}

	fmt.Printf("\n📦 Installing packages: %s\n", strings.Join(manifest.Dependencies.Pacman, ", "))

	if dryRun {
		fmt.Println("   [DRY RUN] Would install packages")
		return nil
	}

	args := append([]string{"pacman", "-S", "--needed", "--noconfirm"}, manifest.Dependencies.Pacman...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func restartServices(manifest *Manifest, dryRun bool) error {
	if len(manifest.Services.Restart) == 0 {
		return nil
	}

	fmt.Println("\n🔄 Restarting services...")

	for _, service := range manifest.Services.Restart {
		if dryRun {
			fmt.Printf("   [DRY RUN] Would restart: %s\n", service)
			continue
		}

		fmt.Printf("   Restarting: %s\n", service)
		cmd := exec.Command("systemctl", "--user", "restart", service)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("restart %s: %w", service, err)
		}
	}
	return nil
}

func restoreBackup(b *Backup) {
	fmt.Println("🔄 Restoring from backup...")
	for destRel, dest := range b.Files {
		backupFile := filepath.Join(b.Path, destRel)
		if _, err := os.Stat(backupFile); os.IsNotExist(err) {
			continue
		}
		os.Remove(dest)
		copyFile(backupFile, dest)
		fmt.Printf("   Restored: %s\n", dest)
	}
	fmt.Printf("✅ Backup restored from: %s\n", b.Path)
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}