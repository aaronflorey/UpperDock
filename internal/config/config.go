package config

import (
	"fmt"
	"github.com/num30/config"
	"log"
	"os"
	"strings"
)

type Config struct {
	Organisation string     `validate:"required"`
	Repository   string     `validate:"required"`
	Token        string     `validate:"required"`
	Interval     int        `validate:"required" default:"900"`
	DockerMode   DockerMode `default:"Docker" validate:"required,oneof=Docker Swarm"`
	GitMode      GitMode    `default:"Release" validate:"required,oneof=Release Package Latest"`
}

func GetConfig() Config {
	var conf Config

	err := config.NewConfReader("config").
		WithSearchDirs("../").
		WithSearchDirs("./").
		Read(&conf)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Println("Running with the following config:")
	log.Printf("ORGANISATION: %s", conf.Organisation)
	log.Printf("REPOSITORY: %s", conf.Repository)
	log.Printf("TOKEN: %s", fmt.Sprintf("%s%s%s", conf.Token[0:7], strings.Repeat("*", 12), conf.Token[len(conf.Token)-4:]))
	log.Printf("INTERVAL: %d", conf.Interval)
	log.Printf("DOCKER_MODE: %s", conf.DockerMode)
	log.Printf("GIT_MODE: %s", conf.GitMode)

	return conf
}
