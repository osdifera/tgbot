package main

import (
	"os"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/joho/godotenv"
	"github.com/oscaletta/chatbot/modules/help"
	"github.com/oscaletta/chatbot/modules/token"
	"github.com/oscaletta/chatbot/modules/welcome"
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

	help.LoadLiquidity(updater)
	welcome.LoadWelcome(updater)
	token.LoadToken(updater)

	logger.Sugar().Info("Updater started successfully")
	updater.StartCleanPolling()
	updater.Idle()
}


