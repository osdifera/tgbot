package help

import (
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
	"github.com/oscaletta/chatbot/config"
)

var markup ext.InlineKeyboardMarkup
var markdownHelpText string

func help(b ext.Bot, u *gotgbot.Update) error {
	msg := b.NewSendableMessage(u.EffectiveChat.Id, "Hey there! I'm Altstrobot."+
		"Commands are preceded with a slash  (/) or an exclamation mark (!)\n\n"+
		"Some basic commands:\n\n"+
		"- /help: for info on how to use me\n\n"+
		"- /liq: get liquidity for a specific pair\n\n"+
		"- /top: get top holders for a contract\n\n\n"+
		"Have fun!")
	msg.ParseMode = parsemode.Html
	msg.ReplyToMessageId = u.EffectiveMessage.MessageId
	msg.ReplyMarkup = &markup
	_, err := msg.Send()
	if err != nil {
		msg.ReplyToMessageId = 0
		_, err = msg.Send()
	}
	return err
}

func markdownHelp(_ ext.Bot, u *gotgbot.Update) error {

	markdownHelpText = "You can use markdown to make your messages more expressive. This is the markdown currently "

	chat := u.EffectiveChat
	if chat.Type != "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in PM!")
		return err
	}

	_, err := u.EffectiveMessage.ReplyHTML(markdownHelpText)
	return err
}

func LoadLiquidity(u *gotgbot.Updater){
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("help", config.BotConfig.Prefix, help))
	u.Dispatcher.AddHandler(handlers.NewCallback("help", markdownHelp))
}