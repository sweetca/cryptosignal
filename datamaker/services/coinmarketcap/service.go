package coinmarketcap

import (
	"fmt"
	"github.com/sweetca/cryptosignal/common/model"
	"go.uber.org/zap"
	"gopkg.in/robfig/cron.v2"
	"sync"
	"time"

	"github.com/sweetca/cryptosignal/common/redis"
	"github.com/sweetca/cryptosignal/datamaker/config"
)

const (
	cronEachMinute = "0 */1 * * * ?"
	serviceName    = "coinmarketcap"
)

type Service struct {
	api        *Api
	zLog       *zap.Logger
	rService   *redis.Service
	cryptoMap  sync.Map
	cryptoList []string
	convert    string
}

func NewService(zLog *zap.Logger, settings *config.Config, rService *redis.Service) (*Service, error) {
	api, err := NewApi(settings)
	if err != nil {
		return nil, fmt.Errorf("fail init coinmarketcap service: %v", err)
	}

	service := Service{
		api:        api,
		zLog:       zLog,
		rService:   rService,
		cryptoList: settings.CryptoList,
		convert:    settings.CoinMarketCapConvert,
	}

	err = service.CacheCryptoMap()
	if err != nil {
		return nil, fmt.Errorf("fail init coin market cap service: %v", err)
	}
	zLog.Debug("coin market cap service: crypto map cached")

	c := cron.New()
	_, err = c.AddFunc(cronEachMinute, service.ProceedCryptoListings)
	if err != nil {
		return nil, err
	}
	c.Start()

	return &service, nil
}

func (s *Service) CacheCryptoMap() error {
	items, err := s.api.Map()
	if err != nil {
		return fmt.Errorf("failt to cache crypto map: %v", err)
	}

	for _, asset := range items {
		s.cryptoMap.Store(asset.Symbol, asset)
	}

	return nil
}

func (s *Service) GetCryptoList() []CryptoItem {
	assetList := make([]CryptoItem, 0)
	s.cryptoMap.Range(func(key, value interface{}) bool {
		assetList = append(assetList, value.(CryptoItem))
		return true
	})
	return assetList
}

func (s *Service) ProceedCryptoListings() {
	items, err := s.api.Listings()
	if err != nil {
		s.zLog.Error("coin market cap service fail to proceed crypto listings", zap.Error(err))
	}

	tNow := time.Now().UTC()
	tNow = tNow.Truncate(1 * time.Minute)
	timestamp := tNow.Unix()

	assetStatistic := make(map[string]CryptoItemStatistic)
	for _, crypto := range items {
		assetStatistic[crypto.Symbol] = crypto
	}

	statList := make([]model.CryptoStatistic, 0)
	for _, symbol := range s.cryptoList {
		stat := assetStatistic[symbol]
		quote := stat.Quote[s.convert]

		modelStat := model.CryptoStatistic{
			Symbol:            stat.Symbol,
			Rank:              stat.CmcRank,
			MaxSupply:         stat.MaxSupply,
			CirculatingSupply: stat.CirculatingSupply,
			TotalSupply:       stat.TotalSupply,
			ConvertedTo:       s.convert,
			Price:             quote.Price,
			Volume24h:         quote.Volume24h,
			MarketCap:         quote.MarketCap,
			LastUpdated:       quote.LastUpdated,
			Timestamp:         timestamp,
		}

		statList = append(statList, modelStat)
	}

	err = s.rService.StoreStatistic(statList)
	if err != nil {
		s.zLog.Error("coin market cap service fail to store crypto stat in redis", zap.Error(err))
	}

	notification := model.ChannelNotification{
		Timestamp: timestamp,
		Action:    redis.ChannelWriteStatDone,
		Source:    serviceName,
		Data:      s.cryptoList,
	}

	err = s.rService.Publish(redis.ChannelWriteStatDone, notification)
	if err != nil {
		s.zLog.Error("coin market cap service fail to notify within redis",
			zap.Any("channel", redis.ChannelWriteStatDone), zap.Error(err))
	}
}
