package github

import (
	"context"
	"dockerup8/internal/config"
	"dockerup8/pkgs/utils"
	"errors"
	"fmt"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"strings"
)

type GHClient struct {
	Context context.Context
	Client  *github.Client
	Config  config.Config
}

func NewGithubClient(config config.Config) *GHClient {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &GHClient{
		Config:  config,
		Context: context.Background(),
		Client:  client,
	}
}

func (g *GHClient) GetLatestRelease() (*string, error) {
	latestRelease, _, err := g.Client.Repositories.GetLatestRelease(g.Context, g.Config.Organisation, g.Config.Repository)

	if err != nil {
		return nil, err
	}

	return latestRelease.TagName, nil
}

func (g *GHClient) GetLatestPackage() (*string, error) {
	packages, err := g.getPackageVersions()

	if err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		return nil, errors.New("no package versions found")
	}

	for _, packageVersion := range packages {
		for _, tag := range packageVersion.Metadata.Container.Tags {
			if strings.HasPrefix(tag, "v") {
				return &tag, nil
			}
		}
	}

	return nil, errors.New("no latest tag found")
}

func (g *GHClient) GetLatestTag() (*string, error) {
	packageVersions, err := g.getPackageVersions()

	if err != nil {
		return nil, err
	}

	if len(packageVersions) == 0 {
		return nil, errors.New("no package versions found")
	}

	for _, packageVersion := range packageVersions {
		if utils.InArray("latest", packageVersion.Metadata.Container.Tags) {
			if packageVersion.Name != nil {
				version := fmt.Sprintf("latest@%s", *packageVersion.Name)
				return &version, nil
			}
		}
	}

	return nil, errors.New("no latest tag found")
}

func (g *GHClient) getPackageVersions() ([]*github.PackageVersion, error) {
	packageType := "container"
	state := "active"

	packageVersions, _, err := g.Client.Organizations.PackageGetAllVersions(
		g.Context,
		g.Config.Organisation,
		packageType,
		g.Config.Repository,
		&github.PackageListOptions{
			Visibility:  nil,
			PackageType: nil,
			State:       &state,
			ListOptions: github.ListOptions{},
		},
	)

	return packageVersions, err
}
