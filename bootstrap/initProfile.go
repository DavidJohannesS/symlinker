package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"symlinker/fs"
	"symlinker/git"
	"symlinker/msg"
)

func InitConfig(path string, repoURL string) {
	absPath := fs.ExpandHome(path)

	if fs.Exists(absPath) {
		msg.Info("Config found at: " + absPath)
		return
	}

	msg.Info("Config missing. Bootstrapping from: " + repoURL)

	parent := filepath.Dir(absPath)
	if err := fs.EnsureDir(parent); err != nil {
		msg.Error(fmt.Sprintf("Failed to create directory %s: %v", parent, err))
		os.Exit(1)
	}

	if err := git.Clone(repoURL, "origin", absPath); err != nil {
		msg.Error(fmt.Sprintf("Failed to clone config repo: %v", err))
		os.Exit(1)
	}
	msg.Info("Bootstrap complete.")
}
