package coinmarketcap

type Response struct {
	Data   interface{} `json:"data"`
	Status Status      `json:"status"`
}

type MapResponse struct {
	Data   []CryptoItem `json:"data"`
	Status Status       `json:"status"`
}

type ListingsResponse struct {
	Data   []CryptoItemStatistic `json:"data"`
	Status Status                `json:"status"`
}

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int32  `json:"elapsed"`
	CreditCount  int32  `json:"credit_count"`
}

type CryptoItem struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	IsActive int    `json:"is_active"`
}

type CryptoItemStatistic struct {
	Id                int32                `json:"id"`
	Symbol            string               `json:"symbol"`
	CmcRank           int                  `json:"cmc_rank"`
	MaxSupply         float32              `json:"max_supply"`
	CirculatingSupply float32              `json:"circulating_supply"`
	TotalSupply       float32              `json:"total_supply"`
	Quote             map[string]QuoteInfo `json:"quote"`
}

type QuoteInfo struct {
	Price       float64 `json:"price"`
	Volume24h   float64 `json:"volume_24h"`
	MarketCap   float64 `json:"market_cap"`
	LastUpdated string  `json:"last_updated"`
}
