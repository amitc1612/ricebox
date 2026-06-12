package rice

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Generator struct {
	OutputPath string
	Name       string
}

func NewGenerator(outputPath, name string) *Generator {
	return &Generator{
		OutputPath: outputPath,
		Name:       name,
	}
}

// commonConfigs maps config file paths to their rice structure names
var commonConfigs = map[string]string{
	".config/hypr/hyprland.conf": "hypr/hyprland.conf",
	".config/waybar/config":      "waybar/config",
	".config/waybar/style.css":   "waybar/style.css",
	".config/kitty/kitty.conf":   "kitty/kitty.conf",
	".config/dunst/dunstrc":      "dunst/dunstrc",
	".config/rofi/config.rasi":   "rofi/config.rasi",
}

func (g *Generator) Generate() error {
	fmt.Printf("🔍 Scanning for configs...\n")

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("can't find home dir: %w", err)
	}

	manifest := &Manifest{
		Name:        g.Name,
		Version:     "0.1.0",
		Description: fmt.Sprintf("Auto-generated rice: %s", g.Name),
		Dotfiles:    make(map[string]string),
	}

	// Create output directories
	dotfilesPath := filepath.Join(g.OutputPath, "dotfiles")
	for _, relPath := range commonConfigs {
		if err := os.MkdirAll(filepath.Join(dotfilesPath, filepath.Dir(relPath)), 0755); err != nil {
			return err
		}
	}

	// Copy existing configs
	for srcRel, destRel := range commonConfigs {
		src := filepath.Join(home, srcRel)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}

		dest := filepath.Join(dotfilesPath, destRel)
		if err := copyFile(src, dest); err != nil {
			return fmt.Errorf("copy %s: %w", srcRel, err)
		}
		manifest.Dotfiles[destRel] = srcRel
		fmt.Printf("   ✓ %s\n", srcRel)
	}

	// Detect installed packages that might be relevant
	manifest.Dependencies.Pacman = detectRelevantPackages()

	// Save manifest
	if err := os.MkdirAll(g.OutputPath, 0755); err != nil {
		return err
	}

	return manifest.Save(g.OutputPath)
}

func detectRelevantPackages() []string {
	relevant := []string{
		"hyprland", "waybar", "kitty", "dunst", "rofi",
		"noto-fonts-emoji", "ttf-jetbrains-mono",
	}

	var installed []string
	for _, pkg := range relevant {
		cmd := exec.Command("pacman", "-Q", pkg)
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