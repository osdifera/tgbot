package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotName       string
	Prefix        []rune
}

var BotConfig Config

// Returns a config object generated from the dotenv file
func init() {
	var err error
	err = godotenv.Load("dev.env")
	if err != nil {
		panic(err.Error())
	}
	returnConfig := Config{}

	// Assign
	var prefixTemp []string
	var runeTemp []rune

	returnConfig.BotName = os.Getenv("BOT_NAME")


	for _, pref := range prefixTemp {
		runeTemp = append(runeTemp, []rune(pref)...)
	}
	returnConfig.Prefix = runeTemp

	// Check Part
	if prefixTemp == nil || len(prefixTemp) < 1 {
		returnConfig.Prefix = []rune{'/', '!'}
		log.Println("[Info][Config] Prefix is not defined, Selecting /")
	}

	BotConfig = returnConfig
}
