package utilities

import (
	"fmt"
	"strconv"

	"github.com/oscaletta/chatbot/models/holder"
	"github.com/oscaletta/chatbot/models/liquidity"
	"github.com/oscaletta/chatbot/models/token"
)

func BuildLiquidityMessage(liquidity liquidity.Response) string {
	var message string = "\n"
	message += " Pair: " + liquidity.Pair.Token0.Symbol +"/"+ liquidity.Pair.Token1.Symbol + "\n"
	message += "Reserva " + liquidity.Pair.Token0.Symbol + ": " + liquidity.Pair.Reserve0 +"\n"
	message += "Reserva " + liquidity.Pair.Token1.Symbol + ": " + liquidity.Pair.Reserve1 +"\n"
	message += "Liquidez total en USD: " + liquidity.Pair.ReserveUSD +"\n"
	return message
}

func BuildTopHoldersMessage(holders []holder.Holder) string {
	var message string = "\n"
	for i,_ := range holders {
		message += "Holder "+strconv.Itoa(i+1) +": "+ holders[i].Address +" \n Porcentaje: "+ fmt.Sprintf("%f", holders[i].Share) +"%"+ "\n"
	}
	return message
}

func BuildArbitrageMessage(token token.Token) string {
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
