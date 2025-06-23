package main

import log "github.com/lukemassa/clilog"

func main() {
	log.SetLogLevel(log.LevelDebug)

	log.Infof("Hello %s", "World")
	log.Debug("initializing subsystems")

	log.Warn("deprecated feature in use")
	log.Errorf("could not read config file")
}
