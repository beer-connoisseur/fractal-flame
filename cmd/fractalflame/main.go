package main

import (
	"log/slog"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/application"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/cli"
)

func main() {
	app := cli.NewApp()

	err := app.Run(os.Args)
	if err != nil {
		slog.Error("Application failed", "error", err)
	}
	exitCode := application.GetExitCode(err)

	os.Exit(exitCode)
}
