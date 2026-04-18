package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"symlinker/core/action"
	"symlinker/core/bootstrap"
	"symlinker/core/config"
	"symlinker/core/fs"
	"symlinker/core/git"
	"symlinker/core/msg"
	"sync"
)
func main() {
	//  Define Flags
	configDirFlag:= flag.String("config-dir", "~/.config/symlinker", "")
	configSource := flag.String("config-repo", "git@github.com:DavidJohannesS/symlinkerConfig.git", "Git URL for bootstrapping config")
	profileNameFlag := flag.String("profile", "desktop", "")
	verbose := flag.Bool("v",false,"")
	dryRun := flag.Bool("dry-run",false,"")
	flag.Parse()

	fmt.Println("Welcome to Dot File Manager reborne")
	fmt.Println("Starting dotFM-reborne!...")
	fmt.Println("---------------------------------")

	bootstrap.InitConfig(*configDirFlag,*configSource)
	configDir := fs.ExpandHome(*configDirFlag)
	mainProfile := filepath.Join(configDir,"config.yaml")
	repoSpecs := filepath.Join(configDir,"repoSpecs")
	cfg, err := config.LoadConfig(
		mainProfile,
		repoSpecs,
	)
	if err != nil {
		panic(err)
	}
	run := &action.Runner{
		IsDryRun: *dryRun,
		IsVerbose: *verbose,
	}
	profileName := *profileNameFlag
	profile, ok := cfg.Profiles[profileName]
	if !ok {
		panic("profile not found: " + profileName)
	}
	home, _ := os.UserHomeDir()
	var wg sync.WaitGroup
	errChan := make(chan string, len(profile.Repos))

	for repoName, repoProfile := range profile.Repos {
		wg.Add(1)
		go func(name string, prof string) {
			defer wg.Done()

			repoCfg, ok := cfg.Repos[name]
			if !ok {
				errChan <- fmt.Sprintf("[%s] Repo configuration not found", name)
				return
			}

			basePath := filepath.Join(home, repoCfg.Path)
			clonePath := filepath.Join(basePath, name)

			run.Do("Ensure directory: "+basePath, func() error {
				return fs.EnsureDir(basePath)
			})

			if fs.Exists(clonePath) {
				run.Do("Pulling repo: "+name, func() error {
					return git.PullRepo(repoCfg.Remote, clonePath)
				})
			} else {
				run.Do("Cloning repo: "+name, func() error {
					return git.Clone(repoCfg.URL, repoCfg.Remote, clonePath)
				})
			}

			repoProfileCfg, ok := repoCfg.Profiles[prof]
			if !ok {
				errChan <- fmt.Sprintf("[%s] Profile '%s' not found", name, prof)
				return
			}

			for _, link := range repoProfileCfg.Links {
				src := filepath.Join(clonePath, link.From)
				dest := fs.ExpandHome(link.To)

				err := run.Do(fmt.Sprintf("[%s] Link %s -> %s", name, link.From, link.To), func() error {
					return fs.CreateSymlink(src, dest)
				})

				if err != nil {
					errChan <- fmt.Sprintf("[%s] Symlink Error (%s): %v", name, link.To, err)
				}
			}
		}(repoName, repoProfile)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for e := range errChan {
		msg.Error(e)
	}
	fmt.Println("---------------------------------")
	fmt.Println("Execution finished.")
}
