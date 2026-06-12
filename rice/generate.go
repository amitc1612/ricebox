package rice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var knownConfigs = map[string]string{
	".config/hypr/hyprland.conf": "hypr/hyprland.conf",
	".config/waybar/config":      "waybar/config",
	".config/waybar/style.css":   "waybar/style.css",
	".config/kitty/kitty.conf":   "kitty/kitty.conf",
	".config/dunst/dunstrc":      "dunst/dunstrc",
	".config/rofi/config.rasi":   "rofi/config.rasi",
	".config/gtk-3.0/settings.ini": "gtk/settings.ini",
	".config/wallpaper":           "wallpaper",
}

func Generate(outputPath, name, description string) error {
	fmt.Printf("🔍 Generating rice: %s\n", name)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("can't find home directory: %w", err)
	}

	manifest := &Manifest{
		Name:        name,
		Version:     "0.1.0",
		Description: description,
		Dotfiles:    make(map[string]string),
	}

	dotfilesDir := filepath.Join(outputPath, "dotfiles")
	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for srcRel, destRel := range knownConfigs {
		src := filepath.Join(home, srcRel)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}

		dest := filepath.Join(dotfilesDir, destRel)
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", destRel, err)
		}

		if err := copyFile(src, dest); err != nil {
			fmt.Printf("   ⚠ skipping %s: %v\n", srcRel, err)
			continue
		}

		manifest.Dotfiles[destRel] = srcRel
		fmt.Printf("   ✓ %s\n", srcRel)
	}

	manifest.Dependencies.Pacman = detectInstalledPackages()

	if err := manifest.Save(outputPath); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Printf("\n✅ Rice generated: %s/\n", outputPath)
	return nil
}

func detectInstalledPackages() []string {
	relevant := []string{
		"hyprland", "waybar", "kitty", "dunst", "rofi",
		"noto-fonts-emoji", "ttf-jetbrains-mono",
		"ttf-nerd-fonts-symbols", "papirus-icon-theme",
	}

	var installed []string
	for _, pkg := range relevant {
		cmd := exec.Command("pacman", "-Qi", pkg)
		if err := cmd.Run(); err == nil {
			installed = append(installed, pkg)
		}
	}
	return installed
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}