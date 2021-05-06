package model

import (
	"encoding/json"
	"fmt"
)

type CryptoStatistic struct {
	Symbol            string  `json:"symbol"`
	Rank              int     `json:"rank"`
	MaxSupply         float32 `json:"max_supply"`
	CirculatingSupply float32 `json:"circulating_supply"`
	TotalSupply       float32 `json:"total_supply"`
	ConvertedTo       string  `json:"converted_to"`

	Price       float64 `json:"price"`
	Volume24h   float64 `json:"volume_24h"`
	MarketCap   float64 `json:"market_cap"`
	LastUpdated string  `json:"last_updated"`

	Timestamp int64 `json:"timestamp"`
}

func (s *CryptoStatistic) GetRedisKey() string {
	return fmt.Sprintf("%s.market.statistic.%d", s.Symbol, s.Timestamp)
}

func (s CryptoStatistic) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("fail CryptoStatistic MarshalBinary: %v : error :%v", s, err)
	}
	return data, nil
}

func (s *CryptoStatistic) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

type ChannelNotification struct {
	Timestamp int64       `json:"timestamp"`
	Action    string      `json:"action"`
	Source    string      `json:"source"`
	Data      interface{} `json:"data"`
}

func (n ChannelNotification) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("fail ChannelNotification MarshalBinary: %v : error :%v", n, err)
	}
	return data, nil
}

func (n *ChannelNotification) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, n)
}
