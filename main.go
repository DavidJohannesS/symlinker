package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"symlinker/config"
	"symlinker/fs"
	"symlinker/git"
	"symlinker/bootstrap"
	"symlinker/msg"
	"sync"
)
func main() {
	//  Define Flags
	configDirFlag:= flag.String("config-dir", "~/.config/symlinker", "")
	dryRun:= flag.Bool("dryRyn", false, "")
	configSource := flag.String("config-repo", "git@github.com:DavidJohannesS/symlinkerConfig.git", "Git URL for bootstrapping config")
	flag.Parse()

	fmt.Println("Welcome to Dot File Manager reborne")
	fmt.Println("Starting dotFM-reborne!...")
	fmt.Println("---------------------------------")
	bootstrap.InitConfig(*configDirFlag,*configSource)
	fmt.Println(*configDirFlag)
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
	
	fmt.Println(cfg)
	if *dryRun {
		os.Exit(100)
	}

	profileName := "desktop"
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

			msg.Info("Processing: " + name)
			basePath := filepath.Join(home, repoCfg.Path)
			clonePath := filepath.Join(basePath, name)

			fs.EnsureDir(basePath)
			if fs.Exists(clonePath) {
				git.PullRepo(repoCfg.Remote, clonePath)
			} else {
				if err := git.Clone(repoCfg.URL, repoCfg.Remote, clonePath); err != nil {
					errChan <- fmt.Sprintf("[%s] Clone Error: %v", name, err)
					return
				}
			}
			repoProfileCfg, ok := repoCfg.Profiles[prof]
			if !ok {
				errChan <- fmt.Sprintf("[%s] Profile '%s' not found in repo config", name, prof)
				return
			}

			for _, link := range repoProfileCfg.Links {
				src := filepath.Join(clonePath, link.From)
				dest := fs.ExpandHome(link.To)
				if err := fs.CreateSymlink(src, dest); err != nil {
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
