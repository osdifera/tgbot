package liquidity

type Response struct {
	Pair Pair `json:"pair"`
}

type Pair struct {
	Reserve0 string `json:"reserve0"`
	Reserve1 string `json:"reserve1"`
	ReserveUSD string `json:"reserveUSD"`
	Token0Price string `json:"token0Price"`
	Token1Price string `json:"token1Price"`
	Token0 Token `json:"token0"`
	Token1 Token `json:"token1"`
}

type Token struct {
	ID string `json:"id"`
	Symbol string `json:"symbol"`
}