package main

import (
    logger "github.com/Chelaran/yagalog"
)

func main() {
    l, err := logger.NewLogger("app.log")
    if err != nil { panic(err) }
    defer l.Close()

    l.WithJSON()
    l.Info("service started: %s", "api")
    l.Error("failed connect: %s", "db")
}


