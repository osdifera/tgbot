package utilities

var TokenList = make(map[string]string)

func LoadTokenList() {
	TokenList["xor"] = "sora"
	TokenList["val"] = "sora-validator-token"
	TokenList["link"] = "chainlink"
	TokenList["shitcoin"] = "shitcoin"
}