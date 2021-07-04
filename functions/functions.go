package functions

import (
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/oscaletta/chatbot/utilities"
	cg "github.com/superoo7/go-gecko/v3"
)

func getTokenId(tokenTicker string) string {
	if _, ok := utilities.TokenList[tokenTicker]; ok {
		return utilities.TokenList[tokenTicker]
	}
	return utilities.TokenList["shitcoin"]
}

func getTokenPrice(tokenName string) string {
	cg := cg.NewClient(nil)
	price, err := cg.SimpleSinglePrice(utilities.TokenList[tokenName[1:]], "usd")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s is worth %f %s", tokenName[1:], price.MarketPrice, price.Currency)
}

func GetTokenPrice(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, getTokenPrice(u.EffectiveMessage.Text))
	return nil
}