package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/oscaletta/chatbot/models/liquidity"
	"github.com/oscaletta/chatbot/models/token"
	"github.com/oscaletta/chatbot/utilities"
	cg "github.com/superoo7/go-gecko/v3"
)

var tokenList = make(map[string]string)

func LoadTokenList() {
	tokenList["xor"] = "sora"
	tokenList["val"] = "sora-validator-token"
	tokenList["link"] = "chainlink"
	tokenList["ramp"] = "ramp"
	tokenList["shitcoin"] = "shitcoin"
}


func getTokenId(tokenTicker string) string {
	if _, ok := tokenList[tokenTicker]; ok {
		return tokenList[tokenTicker]
	}
	return tokenList["shitcoin"]
}

func getTokenPrice(tokenName string) string {
	cg := cg.NewClient(nil)
	price, err := cg.SimpleSinglePrice(tokenList[tokenName[1:]], "usd")
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s is worth %f %s", tokenName[1:], price.MarketPrice, price.Currency)
}

func getTokenArbitrage(tokenName string) string {

	endpointURL := "https://api.coingecko.com/api/v3/coins/" + tokenList[tokenName[3:]] + "/tickers"
	resp, err := http.Get(endpointURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var token token.Token
	json.Unmarshal(bodyBytes, &token)

	return utilities.BuildArbitrageMessage(token)
}

func getLiquidityForPair(pairAddress string) string{
	endpointURL := "http://localhost:5000/v2/pair/" + pairAddress
	resp, err := http.Get(endpointURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var LiqResponse liquidity.Response
	json.Unmarshal(bodyBytes, &LiqResponse)
	return utilities.BuildLiquidityMessage(LiqResponse)
}

func GetTokenPrice(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, getTokenPrice(u.EffectiveMessage.Text))
	return nil
}

func GetArbitrage(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, getTokenArbitrage(u.EffectiveMessage.Text))
	return gotgbot.ContinueGroups{} // will keep executing handlers, even after having been caught by this one.
}

func GetLiquidity(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, getLiquidityForPair(u.EffectiveMessage.Text))
	return gotgbot.ContinueGroups{}
}