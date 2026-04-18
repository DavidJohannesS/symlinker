package git

import (
	"os"
	"path"
	"strings"

	"symlinker/core/msg"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func RepoNameFromURL(url string) string {
    base := path.Base(url)
    return strings.TrimSuffix(base, ".git")
}

func Clone(url, remote, path string) error {
	msg.Success("Cloning," + url + "into" + path)
    auth, err := ssh.NewPublicKeysFromFile("git", os.Getenv("HOME")+"/.ssh/id_rsa", "")
    if err != nil {
        return err
    }
	_, err = gogit.PlainClone(path, false, &gogit.CloneOptions{
        URL:  url,
		RemoteName: remote,
        Auth: auth,
    })
    return err
}
func PullRepo(remote, path string) error {
    repo, err := gogit.PlainOpen(path)
    if err != nil {
        return err
    }

    wt, err := repo.Worktree()
    if err != nil {
        return err
    }

    auth, err := ssh.NewPublicKeysFromFile("git", os.Getenv("HOME")+"/.ssh/id_rsa", "")
    if err != nil {
        return err
    }

    err = wt.Pull(&gogit.PullOptions{
        RemoteName: remote,
        Auth:       auth,
    })

    if err == gogit.NoErrAlreadyUpToDate {
        msg.Info("Already up to date: " + path)
        return nil
    }

    return err
}
