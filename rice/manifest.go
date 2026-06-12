package rice

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Manifest struct {
	Name        string            `toml:"name"`
	Version     string            `toml:"version"`
	Description string            `toml:"description"`
	Author      string            `toml:"author,omitempty"`
	Dotfiles    map[string]string `toml:"dotfiles"`
	Dependencies struct {
		Pacman []string `toml:"pacman,omitempty"`
		Aur    []string `toml:"aur,omitempty"`
	} `toml:"dependencies"`
	Theming struct {
		Gtk    string `toml:"gtk,omitempty"`
		Icons  string `toml:"icons,omitempty"`
		Cursor string `toml:"cursor,omitempty"`
		Font   string `toml:"font,omitempty"`
	} `toml:"theming,omitempty"`
	Services struct {
		Restart []string `toml:"restart,omitempty"`
	} `toml:"services"`
	Hooks struct {
		PreSwitch  string `toml:"pre_switch,omitempty"`
		PostSwitch string `toml:"post_switch,omitempty"`
	} `toml:"hooks"`
}

func LoadManifest(path string) (*Manifest, error) {
	manifestPath := filepath.Join(path, "rice.toml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no rice.toml found in %s", path)
	}

	var m Manifest
	if _, err := toml.DecodeFile(manifestPath, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	if m.Name == "" {
		return nil, fmt.Errorf("rice.toml is missing required field: name")
	}

	return &m, nil
}

func (m *Manifest) Save(path string) error {
	manifestPath := filepath.Join(path, "rice.toml")
	f, err := os.Create(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to create manifest: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(m)
}

func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	if m.Version == "" {
		return fmt.Errorf("version is required")
	}
	for src, dest := range m.Dotfiles {
		if src == "" || dest == "" {
			return fmt.Errorf("dotfile entry has empty key or value: %s -> %s", src, dest)
		}
	}
	return nil
}