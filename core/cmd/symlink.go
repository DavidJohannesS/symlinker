package cmd

import (
    "fmt"
    "path/filepath"
    "symlinker/core/action"
    "symlinker/core/config" // Import your config package
    "symlinker/core/fs"
)

// Change the links slice type here
func LinkProfile(run *action.Runner, repoName string, clonePath string, links []config.LinkConfig) []error {
    var errs []error
    for _, link := range links {
        src := filepath.Join(clonePath, link.From)
        dest := fs.ExpandHome(link.To)

        err := run.Do(fmt.Sprintf("[%s] Link %s", repoName, link.To), func() error {
            return fs.CreateSymlink(src, dest,repoName)
        })
        if err != nil {
            errs = append(errs, err)
        }
    }
    return errs
}
