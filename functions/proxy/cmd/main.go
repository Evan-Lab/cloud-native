package main

import (
	"log/slog"
	"os"

	// Blank-import the function package so the init() runs
	_ "github.com/Evan-Lab/cloud-native/functions/proxy"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

	slog.Info("Starting function host", "url", "http://"+hostname+":"+port, "host", hostname, "port", port)
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		slog.Error("funcframework.StartHostPort", "error", err)
		return
	}
}
