package main

import (
	"github.com/gorward/logger"
)

func main() {
	log := logger.New(logger.Config{
		Level: logger.Error, // posible value logger.All, logger.Error, logger.Warn, ...
		Err:   "error.log",  // file path
		Warn:  "warn.log",
		Debug: "debug.log",
		Info:  "output.log",
	})

	log.Error("Error")
	log.Debug("Debug")
	log.Info("Info")
	log.Warn("Warning")
}
