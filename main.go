package main

import (
	"os"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/joho/godotenv"
	"github.com/oscaletta/chatbot/functions"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	var err error
	err = godotenv.Load("dev.env")
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	log := zap.NewProductionEncoderConfig()
	log.EncodeLevel = zapcore.CapitalLevelEncoder
	log.EncodeTime = zapcore.RFC3339TimeEncoder

	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(log), os.Stdout, zap.InfoLevel))

	updater, err := gotgbot.NewUpdater(logger, os.Getenv("TG_KEY"))
	if err != nil {
		logger.Panic("Updater failed to start")
		return
	}

	functions.LoadTokenList()
	logger.Sugar().Info("Tokens list loaded")

	logger.Sugar().Info("Updater started successfully")
	updater.StartCleanPolling()
	//updater.Dispatcher.AddHandler(handlers.NewCommand("romestime", romesTime))
	//updater.Dispatcher.AddHandler(handlers.NewCommand("price", usdPrice))

	//updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)p", returnTokenPrice))
	//updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)arb", executArbitrage))
	//updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)0x",returnLiquidity))
	updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)top",functions.ReturnTopHolders))
	updater.Idle()
}


