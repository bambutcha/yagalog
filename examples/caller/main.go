package main

import (
    logger "github.com/Chelaran/yagalog"
)

func main() {
    l, err := logger.NewLogger("app.log")
    if err != nil { panic(err) }
    defer l.Close()

    l.WithCaller(true)
    l.Info("with caller example")
}


