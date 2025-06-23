package main

import (
	log "github.com/lukemassa/clilog"
)

func main() {
	log.SetLogLevel(log.LevelDebug)
	log.Debug("starting up")
	log.Debugf("Initialing %s", "subsystems")
	log.Warn("deprecated feature in use")
	log.Error("could not read config file")
}
