package welcome

import (
	"log"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
)

func welcome(bot ext.Bot, u *gotgbot.Update, args []string) error {
	u.EffectiveMessage.ReplyText("Welcome")
	return nil
}

func LoadWelcome(u *gotgbot.Updater) {
	defer log.Println("Loading module welcome")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("welcome", []rune{'!', '/'}, welcome))
}