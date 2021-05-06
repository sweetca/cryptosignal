package coinmarketcap

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sweetca/cryptosignal/datamaker/config"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	reqTimeout = time.Second * 5

	keyInfo        = "key/info"
	cryptoMap      = "cryptocurrency/map"
	cryptoListings = "cryptocurrency/listings/latest?convert=%s"

	authH   = "X-CMC_PRO_API_KEY"
	acceptH = "Accept"
	jsonH   = "applicationJson"
)

type Api struct {
	apiKey     string
	apiURL     string
	convert    string
	httpClient *http.Client
}

func NewApi(settings *config.Config) (*Api, error) {
	client := http.Client{
		Timeout: reqTimeout,
	}

	api := Api{
		apiKey:     settings.CoinMarketCapKey,
		apiURL:     settings.CoinMarketCapAPI,
		convert:    settings.CoinMarketCapConvert,
		httpClient: &client,
	}

	err := api.KeyCheck()
	if err != nil {
		return nil, err
	}

	return &api, nil
}

func (a *Api) KeyCheck() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.apiURL, keyInfo), nil)
	if err != nil {
		return err
	}

	data, status, err := request(a.apiKey, req, a.httpClient)
	if err != nil {
		return fmt.Errorf("fail to execute coinmarketcap request: %v", err)
	}

	var info Response
	err = json.Unmarshal(data, &info)
	if err != nil {
		return fmt.Errorf("fail to unmarshar coinmarketcap key info response: %v", err)
	}

	if status != http.StatusOK {
		return fmt.Errorf("fail coinmarketcap key: %s", info.Status.ErrorMessage)
	}

	return nil
}

func (a *Api) Map() ([]CryptoItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.apiURL, cryptoMap), nil)
	if err != nil {
		return nil, err
	}

	data, status, err := request(a.apiKey, req, a.httpClient)
	if err != nil {
		return nil, fmt.Errorf("fail to execute coinmarketcap request: %v", err)
	}

	var mapResponse MapResponse
	err = json.Unmarshal(data, &mapResponse)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshar coinmarketcap map response: %v", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("fail coinmarketcap key: %s", mapResponse.Status.ErrorMessage)
	}

	return mapResponse.Data, nil
}

func (a *Api) Listings() ([]CryptoItemStatistic, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.apiURL, fmt.Sprintf(cryptoListings, a.convert)), nil)
	if err != nil {
		return nil, err
	}

	data, status, err := request(a.apiKey, req, a.httpClient)
	if err != nil {
		return nil, fmt.Errorf("fail to execute coinmarketcap request: %v", err)
	}

	var listingsResponse ListingsResponse
	err = json.Unmarshal(data, &listingsResponse)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshar coinmarketcap listings response: %v", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("fail coinmarketcap listings request: %s", listingsResponse.Status.ErrorMessage)
	}

	return listingsResponse.Data, nil
}

func (a *Api) OHLCV() error {
	//TODO implement for Startup type of API token
	return errors.New("not implemented")
}

func request(auth string, req *http.Request, client *http.Client) ([]byte, int, error) {
	req.Header.Add(authH, auth)
	req.Header.Add(acceptH, jsonH)

	resp, err := client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	return body, resp.StatusCode, nil
}
