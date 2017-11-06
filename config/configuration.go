// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package config

import (
	"github.com/spf13/viper"
)

// Configuration struct
type Configuration struct {
	Logger    *LoggerConfiguration
	Server    *ServerConfiguration
	Database  *DatabaseConfiguration
	Exchanges *ExchangesConfiguration
}

// NewConfiguration constructor
func NewConfiguration(options *viper.Viper) *Configuration {
	return &Configuration{
		Logger: &LoggerConfiguration{
			Prefix:    options.GetString("logger.prefix"),
			LevelName: options.GetString("logger.level"),
		},
		Server: &ServerConfiguration{
			Port:  options.GetInt("server.port"),
			Debug: options.GetBool("server.debug"),
		},
		Database: &DatabaseConfiguration{
			Path: options.GetString("database.path"),
		},
		Exchanges: &ExchangesConfiguration{
			GDAX: &GDAXConfiguration{
				Key:    options.GetString("exchanges.gdax.key"),
				Secret: options.GetString("exchanges.gdax.secret"),
			},
		},
	}
}
