package internal

import (
	"dockerup8/internal/config"
	"dockerup8/internal/container"
	"dockerup8/internal/github"
	"fmt"
	"log"
)

type UpperDock struct {
	config    config.Config
	container *container.Container
}

func NewUpperDock(cfg config.Config) (*UpperDock, error) {

	c, err := container.New(cfg)

	if err != nil {
		return nil, err
	}

	return &UpperDock{
		config:    cfg,
		container: c,
	}, nil
}

func (u *UpperDock) CheckVersion() {

	latestVersion, err := u.GetLatestVersion()

	if err != nil {
		log.Fatal(err)
		return
	}

	if latestVersion == nil {
		log.Fatal("Latest version is nil")
		return
	}

	image := fmt.Sprintf("ghcr.io/%s/%s:%s", u.config.Organisation, u.config.Repository, *latestVersion)

	log.Printf("Latest version: %s", image)
}

func (u *UpperDock) GetLatestVersion() (*string, error) {
	ghClient := github.NewGithubClient(u.config)

	if u.config.GitMode == config.RELEASE {
		return ghClient.GetLatestRelease()
	}

	if u.config.GitMode == config.PACKAGE {
		return ghClient.GetLatestPackage()
	}

	if u.config.GitMode == config.LATEST {
		return ghClient.GetLatestTag()
	}

	return nil, nil
}
