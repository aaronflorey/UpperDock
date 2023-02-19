package cmd

import (
	"dockerup8/internal"
	"dockerup8/internal/config"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Exec() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	config := config.GetConfig()
	client, err := internal.NewUpperDock(config)

	if err != nil {
		panic(err)
		os.Exit(1)
	}

	interval, err := time.ParseDuration(fmt.Sprintf("%ds", config.Interval))

	if err != nil {
		panic(err)
		os.Exit(1)
	}

	ticker := time.NewTicker(interval)

	client.CheckVersion()

	for {
		select {
		case <-ticker.C:
			go client.CheckVersion()
		case <-c:
			ticker.Stop()
			log.Println("Goodbye :)")
			os.Exit(0)
			return
		}
	}
}
