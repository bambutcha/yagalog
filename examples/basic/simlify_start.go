package main

import logger "github.com/Chelaran/yagalog"

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	log.Info("Logger is start!")
}