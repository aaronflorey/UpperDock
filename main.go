package main

import (
	"log"
	"os"
	"regexp"
	"time"
)

func main() {
	config := parseEnvironment()

	ticker := time.NewTicker(config.Interval)
	log.Printf("S %#v", ticker)

	go checkForNewVersion(config)

	for {
		select {
		case <-ticker.C:
			go checkForNewVersion(config)
		}
	}
}

func checkForNewVersion(config Config) {
	release, err := getLatestVersion(config)

	if err != nil {
		log.Println(err)
		return
	}

	if release == nil {
		log.Println("No new version found")
		return
	}

	log.Printf("New version found: %s", *release)
}

func parseEnvironment() Config {

	repo, present := os.LookupEnv("GITHUB_REPO")

	if !present {
		log.Fatal("GITHUB_TOKEN or GITHUB_REPO not set")
		os.Exit(1)
	}

	match, err := regexp.MatchString("^[a-zA-Z0-9-_]+/[a-zA-Z0-9-_]+$", repo)

	if !match || err != nil {
		log.Fatal("Invalid GITHUB_REPO format")
		os.Exit(1)
	}

	token, present := os.LookupEnv("GITHUB_TOKEN")

	if !present {
		log.Fatal("GITHUB_TOKEN or GITHUB_REPO not set")
		os.Exit(1)
	}

	dockerMode, present := os.LookupEnv("DOCKER_MODE")
	if present == false {
		dockerMode = "Docker"
	}

	if DockerMode(dockerMode).isValid() == false {
		log.Fatal("DOCKER_MODE is not valid, please use Docker or Swarm")
		os.Exit(1)
	}

	gitMode, present := os.LookupEnv("GIT_MODE")
	if present == false {
		gitMode = "Release"
	}

	if GitMode(gitMode).isValid() == false {
		log.Fatal("GIT_MODE is not valid, please use Release, Package or Latest")
		os.Exit(1)
	}

	interval, present := os.LookupEnv("INTERVAL")

	if !present {
		interval = "900"
	}

	match, _ = regexp.MatchString("^[0-9]+$", interval)

	if match == false {
		log.Fatal("Interval is not a valid number")
		os.Exit(1)
	}

	duration, err := time.ParseDuration(interval + "s")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	config := Config{
		Repo:       repo,
		Token:      token,
		DockerMode: DockerMode(dockerMode),
		GitMode:    GitMode(gitMode),
		Interval:   duration,
	}

	log.Println("Starting up with the following config")
	log.Printf("GITHUB_REPO: %s", config.Repo)
	log.Printf("GITHUB_TOKEN: %s", config.Token)
	log.Printf("DOCKER_MODE: %s", config.DockerMode)
	log.Printf("GIT_MODE: %s", config.GitMode)
	log.Printf("INTERVAL: %s", config.Interval)

	return config
}
