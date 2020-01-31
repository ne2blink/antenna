package main

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ne2blink/antenna/pkg/antenna"
	"github.com/ne2blink/antenna/pkg/server"
	"github.com/ne2blink/antenna/pkg/storage"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the telegram bot",
		Args:  cobra.NoArgs,
		RunE:  serve,
	}
}

func serve(cmd *cobra.Command, _ []string) error {
	config, err := getConfig(cmd)
	if err != nil {
		return err
	}

	var log *zap.Logger
	if config.GetBool("debug.enabled") {
		log, err = zap.NewDevelopment()
	} else {
		log, err = zap.NewProduction()
		gin.SetMode(gin.ReleaseMode)
	}
	if err != nil {
		return err
	}
	sugar := log.Sugar()

	bot, err := tgbotapi.NewBotAPI(config.GetString("telegram.token"))
	if err != nil {
		return err
	}

	store, err := storage.New(
		config.GetString("storage.type"),
		config.GetStringMap("storage.options"),
	)
	if err != nil {
		return err
	}
	defer store.Close()

	antenna := antenna.New(bot, store, sugar.With("service", "antenna"))
	if config.GetBool("admin.enabled") {
		antenna.AddAdmin(config.GetStringSlice("admin.usernames")...)
	}
	go antenna.Listen()

	server := server.New(store, antenna, sugar.With("service", "server"))
	return server.Listen(config.GetString("http.addr"))
}
