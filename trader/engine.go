// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package trader

import (
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/euskadi31/cryptotrader/database/entity"
	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/rs/zerolog/log"
)

// RunTickerEvent struct
type RunTickerEvent struct {
	Provider string
}

// SubscribeProductEvent struct
type SubscribeProductEvent struct {
	Provider string
	Products []exchanges.Product
}

// Engine struct
type Engine struct {
	db          *storm.DB
	providers   exchanges.Manager
	tickers     map[string]exchanges.TickerProvider
	runTickerCh chan *RunTickerEvent
	productsCh  chan *SubscribeProductEvent
	doneCh      chan bool
}

// NewEngine trader
func NewEngine(db *storm.DB, providers exchanges.Manager) *Engine {
	return &Engine{
		db:          db,
		providers:   providers,
		tickers:     make(map[string]exchanges.TickerProvider),
		runTickerCh: make(chan *RunTickerEvent),
		productsCh:  make(chan *SubscribeProductEvent),
		doneCh:      make(chan bool),
	}
}

func (e *Engine) tradeBuying(provider string, event *exchanges.TickerEvent, campaign *entity.Campaign) {
	if event.Price < campaign.BuyLimit {
		log.Warn().Msgf("Buying %f %s at %f %s", campaign.Volume, event.Product.From, event.Price, event.Product.To)

		campaign.State = entity.CampaignStateSelling

		if err := e.db.Save(campaign); err != nil {
			log.Error().Err(err).Msg("Save Campaign")
		}
	}
}

func (e *Engine) tradeSelling(provider string, event *exchanges.TickerEvent, campaign *entity.Campaign) {
	if event.Price > campaign.SellLimit {
		log.Warn().Msgf("Selling %f %s at %f %s", campaign.Volume, event.Product.From, event.Price, event.Product.To)

		campaign.State = entity.CampaignStateBuying

		if err := e.db.Save(campaign); err != nil {
			log.Error().Err(err).Msg("Save Campaign")
		}
	}
}

func (e *Engine) trade(provider string, event *exchanges.TickerEvent) {
	query := e.db.Select(
		q.Eq("Provider", provider),
		q.Eq("ProductID", event.Product.String()),
		q.In("State", []entity.CampaignState{
			entity.CampaignStateBuy,
			entity.CampaignStateSell,
		}),
	)

	var campaigns []*entity.Campaign

	if err := query.Find(&campaigns); err != nil && err != storm.ErrNotFound {
		log.Error().Err(err).Msg("Find campaigns")

		return
	}

	for _, campaign := range campaigns {
		campaign.Orders = []*entity.Order{}

		// todo populate order into campaign

		if campaign.IsState(entity.CampaignStateBuy) {
			e.tradeBuying(provider, event, campaign)
		} else {
			e.tradeSelling(provider, event, campaign)
		}
	}

	msg := log.Info().
		Str("side", string(event.Side)).
		Float64("volume", event.Size) /*.
		Time("time", event.Time)*/

	if event.Side == exchanges.SideTypeSell {
		msg.Msgf("%s Price: %f ↘", event.Product.From, event.Price)
	} else {
		msg.Msgf("%s Price: %f ↗", event.Product.From, event.Price)
	}
}

func (e *Engine) processEventChannel() {
	for {
		select {
		case evt := <-e.runTickerCh:
			log.Debug().Msgf("Run %s Ticker", evt.Provider)

			go func() {
				ticker := e.tickers[evt.Provider]

				for event := range ticker.Channel() {
					if event.Time.IsZero() {
						continue
					}

					e.trade(evt.Provider, event)
				}
			}()

		case evt := <-e.productsCh:
			productsList := []string{}

			for _, product := range evt.Products {
				productsList = append(productsList, product.String())
			}

			log.Debug().Msgf("Subscribe to product %s on %s exchange", strings.Join(productsList, ", "), evt.Provider)

			if ticker, ok := e.tickers[evt.Provider]; ok {
				if err := ticker.Subscribe(evt.Products...); err != nil {
					log.Error().Err(err).Msgf("Subscribe to product %v", strings.Join(productsList, ", "))
				}
			}

		case <-e.doneCh:
			return
		}
	}
}

func (e *Engine) initProvider(name string) error {
	// provider already init
	if _, ok := e.tickers[name]; ok {
		return nil
	}

	exchange, err := e.providers.Get(name)
	if err != nil {
		return err
	}

	e.tickers[name] = exchange.Ticker()

	e.runTickerCh <- &RunTickerEvent{
		Provider: name,
	}

	return nil
}

func (e *Engine) subscribeProduct(name string, products []exchanges.Product) error {
	// provider already init
	if _, ok := e.tickers[name]; ok == false {
		if err := e.initProvider(name); err != nil {
			return err
		}
	}

	e.productsCh <- &SubscribeProductEvent{
		Provider: name,
		Products: products,
	}

	return nil
}

// SaveCampaign to engine
func (e *Engine) SaveCampaign(campaign *entity.Campaign) error {
	edit := false

	if campaign.ID > 0 {
		edit = true
	}

	if err := e.db.Save(campaign); err != nil {
		return err
	}

	if edit == false {
		if err := e.subscribeProduct(campaign.Provider, []exchanges.Product{
			exchanges.NewProductFromString(campaign.ProductID),
		}); err != nil {
			return err
		}
	}

	return nil
}

// Start engine
func (e *Engine) Start() error {
	go e.processEventChannel()

	var campaigns []*entity.Campaign

	if err := e.db.All(&campaigns); err != nil {
		return err
	}

	providers := map[string]map[string]exchanges.Product{}

	for _, campaign := range campaigns {
		if _, ok := providers[campaign.Provider]; ok == false {
			providers[campaign.Provider] = map[string]exchanges.Product{}
		}

		if _, ok := providers[campaign.Provider][campaign.ProductID]; ok {
			continue
		}

		providers[campaign.Provider][campaign.ProductID] = exchanges.NewProductFromString(campaign.ProductID)
	}

	for provider, products := range providers {
		productList := []exchanges.Product{}

		for _, product := range products {
			productList = append(productList, product)
		}

		e.subscribeProduct(provider, productList)
	}

	return nil
}

// Stop engine
func (e *Engine) Stop() error {
	e.doneCh <- true

	return nil
}
