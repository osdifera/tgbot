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
	"github.com/oscaletta/chatbot/models/liquidity"
	"github.com/oscaletta/chatbot/models/token"
	"github.com/oscaletta/chatbot/utilities"
)

func LoadToken(u *gotgbot.Updater) {
	log.Println("Loading token module")
	u.Dispatcher.AddHandler(handlers.NewRegex(`^\/\btop\b\s(\w+)`, getTopHolders))
	u.Dispatcher.AddHandler(handlers.NewRegex(`^\/\barb\b`, getPrice))
	u.Dispatcher.AddHandler(handlers.NewRegex(`^\/\bliq\b\s(\w+)`, getLiquidity))
	log.Println("token module loaded")
}

func getPrice(b ext.Bot, u *gotgbot.Update) error {
	log.Println("Getting arbitrage prices")
	endpointURL := "https://api.coingecko.com/api/v3/coins/" + utilities.TokenList[u.EffectiveMessage.Text[5:]] + "/tickers"
	
	resp, err := http.Get(endpointURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var token token.Token
	json.Unmarshal(bodyBytes, &token)

	_, err =b.SendMessage(u.Message.Chat.Id, utilities.BuildArbitrageMessage(token))
	return err
}

func getLiquidity(b ext.Bot, u *gotgbot.Update) error {
	defer log.Println("Getting liquidity")

	endpointURL := "http://localhost:8080/liquidity/" + u.EffectiveMessage.Text[5:]
	resp, err := http.Get(endpointURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var LiqResponse liquidity.Response
	json.Unmarshal(bodyBytes, &LiqResponse)

	_,err = b.SendMessage(u.Message.Chat.Id, utilities.BuildLiquidityMessage(LiqResponse))
	return err
}

func getTopHolders(bot ext.Bot, u *gotgbot.Update) error {
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