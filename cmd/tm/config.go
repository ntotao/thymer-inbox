package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	URL                string   `json:"url"`
	Token              string   `json:"token"`
	GitHubToken        string   `json:"github_token"`
	GitHubRepos        []string `json:"github_repos"`
	ReadwiseToken      string   `json:"readwise_token"`
	GoogleClientID     string   `json:"google_client_id"`
	GoogleClientSecret string   `json:"google_client_secret"`
	GoogleCalendars    []string `json:"google_calendars"`
}

func loadConfig() Config {
	config := Config{
		URL:           os.Getenv("THYMER_URL"),
		Token:         os.Getenv("THYMER_TOKEN"),
		GitHubToken:   os.Getenv("GITHUB_TOKEN"),
		ReadwiseToken: os.Getenv("READWISE_TOKEN"),
	}

	if repos := os.Getenv("GITHUB_REPOS"); repos != "" {
		config.GitHubRepos = parseRepoList(repos)
	}

	// Try config file
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "tm", "config")
	data, err := os.ReadFile(configPath)
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if strings.HasPrefix(line, "url=") && config.URL == "" {
				config.URL = strings.TrimPrefix(line, "url=")
			}
			if strings.HasPrefix(line, "token=") && config.Token == "" {
				config.Token = strings.TrimPrefix(line, "token=")
			}
			if strings.HasPrefix(line, "github_token=") && config.GitHubToken == "" {
				config.GitHubToken = strings.TrimPrefix(line, "github_token=")
			}
			if strings.HasPrefix(line, "github_repos=") && len(config.GitHubRepos) == 0 {
				config.GitHubRepos = parseRepoList(strings.TrimPrefix(line, "github_repos="))
			}
			if strings.HasPrefix(line, "readwise_token=") && config.ReadwiseToken == "" {
				config.ReadwiseToken = strings.TrimPrefix(line, "readwise_token=")
			}
			if strings.HasPrefix(line, "google_client_id=") && config.GoogleClientID == "" {
				config.GoogleClientID = strings.TrimPrefix(line, "google_client_id=")
			}
			if strings.HasPrefix(line, "google_client_secret=") && config.GoogleClientSecret == "" {
				config.GoogleClientSecret = strings.TrimPrefix(line, "google_client_secret=")
			}
			if strings.HasPrefix(line, "google_calendars=") && len(config.GoogleCalendars) == 0 {
				config.GoogleCalendars = parseRepoList(strings.TrimPrefix(line, "google_calendars="))
			}
		}
	}

	return config
}

func saveConfig(config Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".config", "tm")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("url=%s\n", config.URL))
	sb.WriteString(fmt.Sprintf("token=%s\n", config.Token))
	
	if config.GitHubToken != "" {
		sb.WriteString(fmt.Sprintf("github_token=%s\n", config.GitHubToken))
	}
	if len(config.GitHubRepos) > 0 {
		sb.WriteString(fmt.Sprintf("github_repos=%s\n", strings.Join(config.GitHubRepos, ",")))
	}
	
	if config.ReadwiseToken != "" {
		sb.WriteString(fmt.Sprintf("readwise_token=%s\n", config.ReadwiseToken))
	}
	
	if config.GoogleClientID != "" {
		sb.WriteString(fmt.Sprintf("google_client_id=%s\n", config.GoogleClientID))
	}
	if config.GoogleClientSecret != "" {
		sb.WriteString(fmt.Sprintf("google_client_secret=%s\n", config.GoogleClientSecret))
	}
	if len(config.GoogleCalendars) > 0 {
		sb.WriteString(fmt.Sprintf("google_calendars=%s\n", strings.Join(config.GoogleCalendars, ",")))
	}

	return os.WriteFile(filepath.Join(configDir, "config"), []byte(sb.String()), 0600)
}

func parseRepoList(s string) []string {
	var repos []string
	for _, r := range strings.Split(s, ",") {
		r = strings.TrimSpace(r)
		if r != "" {
			repos = append(repos, r)
		}
	}
	return repos
}
