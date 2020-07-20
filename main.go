package main

import (
	"flag"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/victorabarros/challenge-bexs/app/server"
	"github.com/victorabarros/challenge-bexs/internal/config"
	"github.com/victorabarros/challenge-bexs/internal/database"
)

var (
	csvName = flag.String("routes", "./input-file.txt", "travel routes file")
)

func main() {
	flag.Parse() // `go run main.go -h` for help flag

	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Error in load Enviromnts variables.")
	}

	loglvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.WithError(err).Fatalf(
			"Error in set log level %s.", cfg.LogLevel)
	}
	logrus.SetLevel(loglvl)

	rots, err := database.New(*csvName)
	if err != nil {
		panic(err)
	}
	defer rots.File.Close()

	// Up Server
	server.Run(rots, *csvName, strconv.Itoa(cfg.Port))
}
