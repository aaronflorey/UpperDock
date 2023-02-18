# UpperDock

## What is this?
There is no standard way to keep your docker containers up to date, there is [Watchtower](https://containrrr.dev/watchtowe), but that does `docker pull` every X minutes and then replaces the running image which isn't the most efficient.

The plan for UpperDock, is to poll the GitHub API for new releases or package versions, either Semver or new hashes under `latest`, and then update the running container with the new image.

Eg. You're running `ghcr.io/google/search:v1.0.0` and a new version is released as `v2.0.0`, this would see the update in the Release API, wait for the docker image to be tagged with that version, and then update Docker to run that new version.

It can also be used to watch for new versions of latest, and pin the container to that specific sha1.

## How to run?
Go fish
