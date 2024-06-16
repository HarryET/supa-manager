package main

import (
	"log/slog"
	"os"
	"supamanager.io/supa-manager/api"
	"supamanager.io/supa-manager/conf"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Loading config...")
	config, err := conf.LoadConfig(".env")
	if err != nil {
		logger.Error("Failed to load configuration, ensure the required environment variables are set.")
		return
	}
	apiInstance, err := api.CreateApi(logger, config)
	if err != nil {
		logger.Error("Failed to start API state. ", err.Error())
		return
	}

	apiInstance.Router().Run(apiInstance.ListenAddress())
}
