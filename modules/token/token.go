package token

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/oscaletta/chatbot/models/holder"
	"github.com/oscaletta/chatbot/utilities"
)

func topHolders(bot ext.Bot, u *gotgbot.Update) error {
	var holders []holder.Holder
	
	endpointURL := "http://localhost:8080/holders/"+u.EffectiveMessage.Text[5:]
	resp, err := http.Get(endpointURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &holders)

	_,err = bot.SendMessage(u.Message.Chat.Id, utilities.BuildTopHoldersMessage(holders))
	return err
}

func LoadToken(u *gotgbot.Updater) {
	defer log.Println("Loading token module")
	u.Dispatcher.AddHandler(handlers.NewRegex(`^\/\btop\b\s(\w+)`, topHolders))
}