package fs

import (
	"symlinker/core/msg"
	"os"
	"path/filepath"
	"strings"
)

func ExpandHome(p string) string {
    home, err := os.UserHomeDir()
	msg.Fail(err)
    if p == "~" {
        return home
    }
	if strings.HasPrefix(p, "~/"){
		return filepath.Join(home,p[2:])
	}
    return p
}
func EnsureDir(path string) error {
    return os.MkdirAll(path, 0755)
}

func Exists(path string) bool {
    _, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
func SymlinkStatus(src, dest string) (exists bool, correct bool, err error) {
    info, err := os.Lstat(dest)
    if os.IsNotExist(err) {
        return false, false, nil
    }
    if err != nil {
        return false, false, err
    }
    if info.Mode()&os.ModeSymlink == 0 {
        return true, false, nil
    }
    target, err := os.Readlink(dest)
    if err != nil {
        return true, false, err
    }

    return true, target == src, nil
}
func CreateSymlink(src string, dest string,repoName string) error {
    exists, correct, err := SymlinkStatus(src, dest)
    if err != nil {
        return err
    }
    if exists && correct {
        msg.Info("["+repoName+"] "+"Symlink already exists: " + dest)
        return nil
    }
    os.Remove(dest)
    return os.Symlink(src, dest)
}
