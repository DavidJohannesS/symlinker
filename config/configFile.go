package config

import (
    "io/fs"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type Config struct {
    Profiles map[string]ProfileConfig `yaml:"profiles"`
    Repos    map[string]RepoConfig    `yaml:"-"`
}

type ProfileConfig struct {
    Repos map[string]string `yaml:"repos"`
}

type RepoConfig struct {
    URL      string                       `yaml:"url"`
    Remote   string                       `yaml:"remote"`
    Path     string                       `yaml:"path"`
    Profiles map[string]RepoProfileConfig `yaml:"profiles"`
}

type RepoProfileConfig struct {
    Links []LinkConfig `yaml:"links"`
}

type LinkConfig struct {
    From string `yaml:"from"`
    To   string `yaml:"to"`
}

func LoadConfig(configPath string, repoSpecsDir string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    repos := make(map[string]RepoConfig)

    err = filepath.WalkDir(repoSpecsDir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }

        if filepath.Ext(path) != ".yaml" {
            return nil
        }

        raw, err := os.ReadFile(path)
        if err != nil {
            return err
        }

        var repo RepoConfig
        if err := yaml.Unmarshal(raw, &repo); err != nil {
            return err
        }

        name := filepath.Base(path)
        name = name[:len(name)-len(filepath.Ext(name))]

        repos[name] = repo
        return nil
    })

    if err != nil {
        return nil, err
    }

    cfg.Repos = repos
    return &cfg, nil
}
