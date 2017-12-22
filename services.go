// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package cryptotrader

import (
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
	"github.com/euskadi31/cryptotrader/services"
	"github.com/euskadi31/cryptotrader/trader"
	"github.com/euskadi31/cryptotrader/trader/algorithms"
	"github.com/euskadi31/go-eventemitter"
	"github.com/euskadi31/go-server"
	"github.com/euskadi31/go-service"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const applicationName = "cryptotrader"

// Service Container
var container = service.New()

// const of service name
const (
	ServiceLoggerKey           string = "service.logger"
	ServiceConfigKey                  = "service.config"
	ServiceRouterKey                  = "service.router"
	ServiceDBKey                      = "service.db.storm"
	ServiceExchangeManagerKey         = "service.exchange.manager"
	ServiceGDAXExchangeKey            = "service.exchange.gdax"
	ServiceTimeseriesKey              = "service.timeseries"
	ServiceTraderEngineKey            = "service.trader.engine"
	ServiceAlgorithmManagerKey        = "service.algorithm.manager"
	ServiceAlgorithmTrendKey          = "service.algorithm.trend"
	ServiceCampaignKey                = "service.campaign"
	ServiceOrderKey                   = "service.order"
	ServiceEventEmitterKey            = "service.eventemitter"
)

func init() {
	container.Set(ServiceLoggerKey, func(c *service.Container) interface{} {
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

	container.Set(ServiceConfigKey, func(c *service.Container) interface{} {
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
		options.SetDefault("logger.prefix", applicationName)
		options.SetDefault("database.path", "/var/lib/cryptotrader")

		options.SetConfigName("config") // name of config file (without extension)

		options.AddConfigPath("/etc/" + applicationName + "/")   // path to look for the config file in
		options.AddConfigPath("$HOME/." + applicationName + "/") // call multiple times to add many search paths
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

	container.Set(ServiceEventEmitterKey, func(c *service.Container) interface{} {
		return eventemitter.New()
	})

	container.Set(ServiceDBKey, func(c *service.Container) interface{} {
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		path := strings.TrimRight(cfg.Database.Path, "/")

		db, err := storm.Open(fmt.Sprintf("%s/cryptotrader.db", path))
		if err != nil {
			log.Fatal().Err(err).Msg(ServiceDBKey)
		}

		return db
	})

	container.Set(ServiceGDAXExchangeKey, func(c *service.Container) interface{} {
		// cfg := c.Get(ServiceConfigKey).(*config.Configuration)

		ex, err := gdax.NewGDAX()
		if err != nil {
			log.Fatal().Err(err).Msg(ServiceGDAXExchangeKey)
		}

		return ex
	})

	container.Set(ServiceExchangeManagerKey, func(c *service.Container) interface{} {
		manager := exchanges.NewManager()

		manager.Add(c.Get(ServiceGDAXExchangeKey).(exchanges.ExchangeProvider))

		return manager
	})

	container.Set(ServiceAlgorithmTrendKey, func(c *service.Container) interface{} {
		campaignService := c.Get(ServiceCampaignKey).(*services.CampaignService)
		orderService := c.Get(ServiceOrderKey).(*services.OrderService)

		return algorithms.NewTrend(campaignService, orderService)
	})

	container.Set(ServiceAlgorithmManagerKey, func(c *service.Container) interface{} {
		manager := algorithms.NewManager()

		manager.Add(c.Get(ServiceAlgorithmTrendKey).(algorithms.Algorithm))

		return manager
	})

	container.Set(ServiceCampaignKey, func(c *service.Container) interface{} {
		db := c.Get(ServiceDBKey).(*storm.DB)

		return services.NewCampaignService(db)
	})

	container.Set(ServiceOrderKey, func(c *service.Container) interface{} {
		db := c.Get(ServiceDBKey).(*storm.DB)

		return services.NewOrderService(db)
	})

	container.Set(ServiceTraderEngineKey, func(c *service.Container) interface{} {
		db := c.Get(ServiceDBKey).(*storm.DB)
		exchangesManager := c.Get(ServiceExchangeManagerKey).(exchanges.Manager)
		algorithmsManager := c.Get(ServiceAlgorithmManagerKey).(algorithms.Manager)
		emitter := c.Get(ServiceEventEmitterKey).(eventemitter.EventEmitter)

		return trader.NewEngine(db, exchangesManager, algorithmsManager, emitter)
	})

	container.Set(ServiceRouterKey, func(c *service.Container) interface{} {
		logger := c.Get(ServiceLoggerKey).(zerolog.Logger)
		cfg := c.Get(ServiceConfigKey).(*config.Configuration)
		db := c.Get(ServiceDBKey).(*storm.DB)
		engine := c.Get(ServiceTraderEngineKey).(*trader.Engine)
		algorithmsManager := c.Get(ServiceAlgorithmManagerKey).(algorithms.Manager)
		emitter := c.Get(ServiceEventEmitterKey).(eventemitter.EventEmitter)

		router := server.NewRouter()

		router.Use(hlog.NewHandler(logger))
		router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			rlog := hlog.FromRequest(r)

			var evt *zerolog.Event

			switch {
			case status >= 200 && status <= 299:
				evt = rlog.Info()
			case status >= 300 && status <= 399:
				evt = rlog.Info()
			case status >= 400 && status <= 499:
				evt = rlog.Warn()
			default:
				evt = rlog.Error()
			}

			evt.
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msgf("%s %s", r.Method, r.URL.Path)
		}))
		router.Use(hlog.RemoteAddrHandler("ip"))
		router.Use(hlog.UserAgentHandler("user_agent"))
		router.Use(hlog.RefererHandler("referer"))
		router.Use(hlog.RequestIDHandler("req_id", "Request-Id"))

		router.EnableRecovery()
		router.EnableHealthCheck()
		router.EnableCorsWithOptions(cors.Options{
			AllowCredentials: true,
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
				"X-Requested-With",
			},
			Debug: cfg.Server.Debug,
		})

		router.SetNotFoundFunc(func(w http.ResponseWriter, r *http.Request) {
			server.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error": map[string]interface{}{
					"message": fmt.Sprintf("%s %s not found", r.Method, r.URL.Path),
				},
			})
		})

		router.AddController(controllers.NewTimeseriesController(engine, emitter))
		router.AddController(controllers.NewCampaignController(db, engine))
		router.AddController(controllers.NewAlgorithmController(algorithmsManager))
		router.AddController(controllers.NewUIController())

		return router
	})
}

// Run Application
func Run() {
	_ = container.Get(ServiceLoggerKey).(zerolog.Logger)
	cfg := container.Get(ServiceConfigKey).(*config.Configuration)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	emitter := container.Get(ServiceEventEmitterKey).(eventemitter.EventEmitter)
	router := container.Get(ServiceRouterKey).(*server.Router)

	engine := container.Get(ServiceTraderEngineKey).(*trader.Engine)

	go func() {
		log.Info().Msg("Starting trader engine...")

		engine.Start()
	}()

	log.Info().Msgf("Server running on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}

	emitter.Wait()

	log.Info().Msg("Server exited")
}
