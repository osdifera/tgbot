package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strconv"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/joho/godotenv"
	"github.com/oscaletta/chatbot/models/holder"
	"github.com/oscaletta/chatbot/models/liquidity"
	"github.com/oscaletta/chatbot/models/token"
	cg "github.com/superoo7/go-gecko/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var tokenList = make(map[string]string)

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

	loadTokenList()
	logger.Sugar().Info("Tokens list loaded")

	logger.Sugar().Info("Updater started successfully")
	updater.StartCleanPolling()
	//updater.Dispatcher.AddHandler(handlers.NewCommand("romestime", romesTime))
	//updater.Dispatcher.AddHandler(handlers.NewCommand("price", usdPrice))

	//updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)p", returnTokenPrice))
	//updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)arb", executArbitrage))
	//updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)0x",returnLiquidity))
	updater.Dispatcher.AddHandler(handlers.NewRegex("(?i)top",returnTopHolders))
	updater.Idle()
}

func loadTokenList() {
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
	//fmt.Printf("API Response as struct %+v\n", token)

	return buildArbitrageMessage(token)
}

func buildArbitrageMessage(token token.Token) string {
	var message string = "\n"
	counter := 1
	for _, ticker := range token.Tickers {
		//for _, market := range ticker
		message += " " + ticker.Market.Name + "\n"
		for _, value := range ticker.ConvertedLast {
			if counter == 1 {
				message += "" + fmt.Sprintf("BTC: %f", value) + "\n"
			}
			if counter == 2 {
				message += "" + fmt.Sprintf("ETH: %f", value) + "\n"
			}
			if counter == 3 {
				message += "" + fmt.Sprintf("USD: %f", value) + "\n"
			}
			counter++
		}
		counter = 1
	}
	return message
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
	return buildLiquidityMessage(LiqResponse)
}

func getTopHolders(address string) string{
	var holders []holder.Holder
	
	endpointURL := "http://localhost:8080/holders/"+address[3:]
	resp, err := http.Get(endpointURL)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bodyBytes, &holders)
	fmt.Println(holders)
	return buildTopHoldersMessage(holders)
}

func buildLiquidityMessage(liquidity liquidity.Response) string {
	var message string = "\n"
	message += " Pair: " + liquidity.Pair.Token0.Symbol +"/"+ liquidity.Pair.Token1.Symbol + "\n"
	message += "Reserva " + liquidity.Pair.Token0.Symbol + ": " + liquidity.Pair.Reserve0 +"\n"
	message += "Reserva " + liquidity.Pair.Token1.Symbol + ": " + liquidity.Pair.Reserve1 +"\n"
	message += "Liquidez total en USD: " + liquidity.Pair.ReserveUSD +"\n"
	return message
}

func buildTopHoldersMessage(holders []holder.Holder) string {
	var message string = "\n"
	for i,_ := range holders {
		message += "Holder "+strconv.Itoa(i+1) +": "+ holders[i].Address +" \n Porcentaje: "+ fmt.Sprintf("%f", holders[i].Share) +"%"+ "\n"
	}
	return message
}

func returnTokenPrice(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, getTokenPrice(u.EffectiveMessage.Text))
	return nil
}

func executArbitrage(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, getTokenArbitrage(u.EffectiveMessage.Text))
	return gotgbot.ContinueGroups{} // will keep executing handlers, even after having been caught by this one.
}

func returnLiquidity(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, getLiquidityForPair(u.EffectiveMessage.Text))
	return gotgbot.ContinueGroups{}
}

func returnTopHolders(b ext.Bot, u *gotgbot.Update) error {
	b.SendMessage(u.Message.Chat.Id, getTopHolders(u.EffectiveMessage.Text))
	return gotgbot.ContinueGroups{}
}