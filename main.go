// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"

	"github.com/euskadi31/cryptotrader/config"
	"github.com/euskadi31/cryptotrader/trader"
	"github.com/euskadi31/go-server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const applicationName = "cryptotrader"

func main() {
	_ = container.Get(ServiceLoggerKey).(zerolog.Logger)
	cfg := container.Get(ServiceConfigKey).(*config.Configuration)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	router := container.Get(ServiceRouterKey).(*server.Router)

	engine := container.Get(ServiceTraderEngineKey).(*trader.Engine)

	go func() {
		log.Info().Msg("Starting trader engine...")

		engine.Start()
	}()

	log.Info().Msgf("Server running on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	log.Info().Msg("Server exited")
}
