package config

type DockerMode string
type GitMode string

const (
	DOCKER  DockerMode = "Docker"
	SWARM   DockerMode = "Swarm"
	RELEASE GitMode    = "Release"
	PACKAGE GitMode    = "Package"
	LATEST  GitMode    = "Latest"
)

func (d DockerMode) isValid() bool {
	switch d {
	case DOCKER, SWARM:
		return true
	}

	return false
}

func (g GitMode) isValid() bool {
	switch g {
	case RELEASE, PACKAGE, LATEST:
		return true
	}

	return false
}
