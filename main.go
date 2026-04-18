package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"symlinker/core/action"
	"symlinker/core/bootstrap"
	"symlinker/core/cmd"
	"symlinker/core/config"
	"symlinker/core/fs"
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

			repoCfg, _ := cfg.Repos[name]
			basePath := filepath.Join(home, repoCfg.Path)

			clonePath, err := cmd.SyncRepo(run, name, repoCfg.Remote, repoCfg.URL, basePath)
			if err != nil {
				errChan <- fmt.Sprintf("[%s] Sync Error: %v", name, err)
				return
			}

			repoProfileCfg, _ := repoCfg.Profiles[prof]
			errs := cmd.LinkProfile(run, name, clonePath, repoProfileCfg.Links)

			for _, e := range errs {
				errChan <- fmt.Sprintf("[%s] Link Error: %v", name, e)
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
