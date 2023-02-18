package main

import (
	"context"
	"errors"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"strings"
)

var repoPackage *github.Package
var ctx = context.Background()

func createGithubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

func getLatestVersion(config Config) (*string, error) {
	client := createGithubClient(ctx, config.Token)
	organisation, repo := getRepo(config)

	if config.GitMode == RELEASE {
		return getLatestRelease(ctx, client, organisation, repo)
	}

	if config.GitMode == PACKAGE {
		return getLatestPackage(ctx, client, organisation, repo, nil)
	}

	if config.GitMode == LATEST {
		tag := "latest"
		return getLatestPackage(ctx, client, organisation, repo, &tag)
	}

	return nil, nil
}

func getLatestPackage(ctx context.Context, client *github.Client, organisation string, repo string, tag *string) (*string, error) {

	packageVersions, err := getPackageVersions(ctx, client, organisation, repo)

	if err != nil {
		return nil, err
	}

	if len(packageVersions) == 0 {
		return nil, errors.New("no package versions found")
	}

	if tag != nil {
		for _, packageVersion := range packageVersions {
			if *packageVersion.Version == *tag {
				return packageVersion.Version, nil
			}
		}

		return nil, errors.New("no package versions found")
	}

	return packageVersions[0].Version, nil
}

func getLatestRelease(ctx context.Context, client *github.Client, organisation string, repo string) (*string, error) {
	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, organisation, repo)

	if err != nil {
		return nil, err
	}

	return latestRelease.TagName, nil
}

func getRepo(config Config) (string, string) {
	res := strings.Split(config.Repo, "/")

	return res[0], res[1]
}

func getPackageVersions(ctx context.Context, client *github.Client, organisation string, repo string) ([]*github.PackageVersion, error) {
	packageType := "container"
	state := "active"

	packageVersions, _, err := client.Organizations.PackageGetAllVersions(
		ctx,
		organisation,
		packageType,
		repo,
		&github.PackageListOptions{
			Visibility:  nil,
			PackageType: nil,
			State:       &state,
			ListOptions: github.ListOptions{},
		},
	)

	return packageVersions, err
}
