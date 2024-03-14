package main

import (
	"context"
	"os"

	"github.com/BradMyrick/goFlat/cmd"
	"github.com/BradMyrick/goFlat/obs"
	"go.uber.org/zap"
)

func main() {
    defer func() {
        if r := recover(); r != nil {
            // Report error appropriately
            obs.Logger(context.Background()).Error("Global Exception Caught", zap.Error(r.(error)))
        }
    }()

    if err := cmd.Execute(); err != nil {
        obs.Logger(context.Background()).Error("Command execution failed", zap.Error(err))
		os.Exit(1)
    }
}
                                                                                                                         