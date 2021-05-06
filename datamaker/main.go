package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"

	"github.com/sweetca/cryptosignal/common/log"
	"github.com/sweetca/cryptosignal/common/redis"
	"github.com/sweetca/cryptosignal/datamaker/config"
	"github.com/sweetca/cryptosignal/datamaker/services/coinmarketcap"
)

var zLog *zap.Logger
var cmcService *coinmarketcap.Service

func main() {
	var err error
	zLog, err = log.NewLog(zapcore.DebugLevel)
	if err != nil {
		panic(err)
	}
	settings := config.Init()
	zLog.Info("data maker config", zap.Any("settings", settings))

	rConfig := redis.Config{
		RedisPort:     settings.RedisPort,
		RedisHost:     settings.RedisHost,
		RedisPoolSize: settings.RedisPoolSize,
	}
	rService, err := redis.NewService(&rConfig)
	if err != nil {
		zLog.Fatal("fail data maker on redis init", zap.Any("error", err))
	}

	cmcService, err = coinmarketcap.NewService(zLog, settings, rService)
	if err != nil {
		zLog.Fatal("fail data maker", zap.Any("error", err))
	}

	port := settings.AppPort
	router := mux.NewRouter()
	router.Use(middleware)

	router.HandleFunc("/info", handleInfoRequest).Methods("GET")

	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET"})
	maxAge := handlers.MaxAge(1000)
	corsHandler := handlers.CORS(originsOk, headersOk, methodsOk, maxAge)(router)

	err = http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), corsHandler)
	if err != nil {
		zLog.Fatal(err.Error())
	}
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zLog.Debug("rest api middleware requested")
		next.ServeHTTP(w, r)
	})
}

func handleInfoRequest(w http.ResponseWriter, _ *http.Request) {
	assetList := cmcService.GetCryptoList()

	var result []byte
	result, _ = json.Marshal(assetList)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, _ = w.Write(result)
}
