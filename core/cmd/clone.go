package cmd

import (
	"path/filepath"
	"symlinker/core/action"
	"symlinker/core/fs"
	"symlinker/core/git"
)

func SyncRepo(run *action.Runner, name string, remote string, url string, basePath string) (string, error) {
	clonePath := filepath.Join(basePath, name)

	if fs.Exists(clonePath) {
		err := run.Do("["+ name + "] " + "Pulling: "+name, func() error {
			return git.PullRepo(remote, clonePath,name)
		})
		return clonePath, err
	}

	err := run.Do("["+name + "] "+"Cloning: "+name, func() error {
		return git.Clone(url, remote, clonePath)
	})
	return clonePath, err
}
