package main

import (
	"github.com/bitrise-io/go-utils/command"
	log "github.com/sirupsen/logrus"
)

func EnvmanInitAtPath(envstorePth string) error {
	logLevel := log.GetLevel().String()
	args := []string{"--loglevel", logLevel, "--path", envstorePth, "init", "--clear"}
	return command.RunCommand("envman", args...)
}
