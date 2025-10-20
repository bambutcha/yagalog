package main

import (
    "log"
    "time"
    logger "github.com/Chelaran/yagalog"
)

func main() {
    l, err := logger.NewLogger("app.log")
	if err != nil { log.Fatal(err) }
	defer l.Close()

    l.SetLevel(logger.INFO)
    l.WithTimeFormat(time.RFC3339Nano)
    l.Info("Hello from YagaLog %d", 1)
	l.Warning("Be careful")
	l.Error("Oops: %s", "something went wrong")
}