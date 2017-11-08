// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/euskadi31/cryptotrader/config"
	"github.com/euskadi31/cryptotrader/controllers"
	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/euskadi31/cryptotrader/exchanges/gdax"
	"github.com/euskadi31/cryptotrader/timeseries"
	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-service"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Service Container
var Service = service.New()

// const of service name
const (
	ServiceLoggerKey       string = "service.logger"
	ServiceConfigKey              = "service.config"
	ServiceRouterKey              = "service.router"
	ServiceDBKey                  = "service.db.storm"
	ServiceGDAXExchangeKey        = "service.exchange.gdax"
	ServiceTimeseriesKey          = "service.timeseries"
)

func init() {
	Service.Set(ServiceLoggerKey, func(c *service.Container) interface{} {
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		logger := zerolog.New(os.Stdout).With().
			Timestamp().
			Str("role", cfg.Logger.Prefix).
			//Str("host", host).
			Logger()

		zerolog.SetGlobalLevel(cfg.Logger.Level())

		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		stdlog.SetFlags(0)
		stdlog.SetOutput(logger)

		log.Logger = logger

		return logger
	})

	Service.Set(ServiceConfigKey, func(c *service.Container) interface{} {
		var cfgFile string
		cmd := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		cmd.StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.yaml)")

		// Ignore errors; cmd is set for ExitOnError.
		cmd.Parse(os.Args[1:])

		options := viper.New()

		if cfgFile != "" { // enable ability to specify config file via flag
			options.SetConfigFile(cfgFile)
		}

		options.SetDefault("server.port", 8989)
		options.SetDefault("logger.level", "info")
		options.SetDefault("logger.prefix", ApplicationName)
		options.SetDefault("database.path", "/var/lib/cryptotrader")

		options.SetConfigName("config") // name of config file (without extension)

		options.AddConfigPath("/etc/" + ApplicationName + "/")   // path to look for the config file in
		options.AddConfigPath("$HOME/." + ApplicationName + "/") // call multiple times to add many search paths
		options.AddConfigPath(".")

		if port := os.Getenv("PORT"); port != "" {
			os.Setenv("CRYPTOTRADER_SERVER_PORT", port)
		}

		options.SetEnvPrefix("CRYPTOTRADER")
		options.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		options.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := options.ReadInConfig(); err == nil {
			log.Info().Msgf("Using config file: %s", options.ConfigFileUsed())
		}

		return config.NewConfiguration(options)
	})

	Service.Set(ServiceDBKey, func(c *service.Container) interface{} {
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		path := strings.TrimRight(cfg.Database.Path, "/")

		db, err := storm.Open(fmt.Sprintf("%s/cryptotrader.db", path))
		if err != nil {
			log.Fatal().Err(err)
		}

		return db
	})

	Service.Set(ServiceGDAXExchangeKey, func(c *service.Container) interface{} {
		// cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		ex, err := gdax.NewGDAX()
		if err != nil {
			log.Fatal().Err(err)
		}

		return ex
	})

	Service.Set(ServiceTimeseriesKey, func(c *service.Container) interface{} {
		return timeseries.New()
	})

	Service.Set(ServiceRouterKey, func(c *service.Container) interface{} {
		logger := c.Get(ServiceLoggerKey).(zerolog.Logger)
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)
		ts := c.Get(ServiceTimeseriesKey).(*timeseries.Timeseries)

		corsHandler := cors.New(cors.Options{
			AllowCredentials: false,
			AllowedOrigins:   []string{"*"},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodOptions,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
			},
			Debug: cfg.Server.Debug,
		})

		router := server.NewRouter()

		router.Use(hlog.NewHandler(logger))
		router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}))
		router.Use(hlog.RemoteAddrHandler("ip"))
		router.Use(hlog.UserAgentHandler("user_agent"))
		router.Use(hlog.RefererHandler("referer"))
		router.Use(hlog.RequestIDHandler("req_id", "Request-Id"))
		router.Use(corsHandler.Handler)

		router.EnableHealthCheck()

		router.AddController(controllers.NewTimeseriesController(ts))

		ex := Service.Get(ServiceGDAXExchangeKey).(exchanges.ExchangeProvider)

		go func() {
			file, err := os.OpenFile("history.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
			defer file.Close()

			writer := csv.NewWriter(file)
			defer writer.Flush()

			log.Info().Msg("Call Ticker")
			c, err := ex.Ticker("BTC", "EUR")
			if err != nil {
				log.Error().Err(err).Msg("")
			}

			for ticker := range c {
				if ticker.Time.IsZero() {
					continue
				}

				ts.Add(ticker.Time.Unix(), ticker.Price)

				if err := writer.Write([]string{
					fmt.Sprintf("%d", ticker.Time.Unix()),
					fmt.Sprintf("%f", ticker.Price),
				}); err != nil {
					log.Error().Err(err).Msg("CSV write failed")
				}

				writer.Flush()

				msg := log.Info().
					Str("side", string(ticker.Side)).
					Float64("volume", ticker.Size).
					Time("time", ticker.Time)

				if ticker.Side == exchanges.SideTypeSell {
					msg.Msgf("BTC Price: %f ↘", ticker.Price)
				} else {
					msg.Msgf("BTC Price: %f ↗", ticker.Price)
				}
			}
		}()

		return router
	})
}

// Run Application
func Run() {
	_ = Service.Get(ServiceLoggerKey).(zerolog.Logger)
	cfg := Service.Get(ServiceConfigKey).(*config.Configuration)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	router := Service.Get(ServiceRouterKey).(*server.Router)

	log.Info().Msgf("Server running on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	log.Info().Msg("Server exited")
}
