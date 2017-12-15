// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package trader

import (
	"fmt"
	"strings"

	"github.com/euskadi31/cryptotrader/trader/algorithms"

	"github.com/euskadi31/cryptotrader/timeseries"

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
	algorithms  algorithms.Manager
	tickers     map[string]exchanges.TickerProvider
	timeseries  map[string]*timeseries.Timeseries
	runTickerCh chan *RunTickerEvent
	productsCh  chan *SubscribeProductEvent
	doneCh      chan bool
}

// NewEngine trader
func NewEngine(db *storm.DB, providers exchanges.Manager, algorithms algorithms.Manager) *Engine {
	return &Engine{
		db:          db,
		providers:   providers,
		algorithms:  algorithms,
		tickers:     make(map[string]exchanges.TickerProvider),
		timeseries:  make(map[string]*timeseries.Timeseries),
		runTickerCh: make(chan *RunTickerEvent),
		productsCh:  make(chan *SubscribeProductEvent),
		doneCh:      make(chan bool),
	}
}

/*
func (e *Engine) tradeBuying(provider string, event *exchanges.TickerEvent, campaign *entity.Campaign, ts *timeseries.Timeseries) {
	if event.Price < campaign.BuyLimit {
		campaign.State = entity.CampaignStateBuying

		if err := e.db.Save(campaign); err != nil {
			log.Error().Err(err).Msg("Save Campaign")
		}

		log.Warn().Msgf("Buying %f %s at %f %s", campaign.Volume, event.Product.From, event.Price, event.Product.To)

		// emulate buying start ----
		order := &entity.Order{
			Provider:  campaign.Provider,
			Side:      exchanges.SideTypeBuy,
			ProductID: campaign.ProductID,
			Size:      campaign.Volume,
			Price:     campaign.Volume * event.Price,
		}

		if err := e.db.Save(order); err != nil {
			log.Error().Err(err).Msg("save order failed")

			return
		}

		campaign.State = entity.CampaignStateSell

		campaign.BuyOrder = order
		// campaign.AddOrder(order)

		if err := e.db.Save(campaign); err != nil {
			log.Error().Err(err).Msg("Save Campaign")
		}
		// end emulate
	}
}

func (e *Engine) tradeSelling(provider string, event *exchanges.TickerEvent, campaign *entity.Campaign, ts *timeseries.Timeseries) {
	switch campaign.SellLimitUnit {
	case "percent":
		log.Debug().Msgf("Current Price: %v", event.Price)
		log.Debug().Msgf("Margin: %v", campaign.BuyOrder.GetMarginInPercent(event.Price))

		if campaign.BuyOrder.GetMarginInPercent(event.Price) < campaign.SellLimit {
			return
		}

		if ts.Size() < 150 {
			log.Debug().Msg("there are not enough elements in the time series")

			return
		}

		a, err := ts.GetTrending(150)
		if err != nil {
			log.Error().Err(err).Msg("GetTrending failed")

			return
		}

		b, err := ts.GetTrending(10)
		if err != nil {
			log.Error().Err(err).Msg("GetTrending failed")

			return
		}

		if a != timeseries.TrendTypeIncreasing && b != timeseries.TrendTypeDecreasing {
			log.Debug().Msg("Not match trend model")

			return
		}

		campaign.State = entity.CampaignStateSelling

		if err := e.db.Save(campaign); err != nil {
			log.Error().Err(err).Msg("Save Campaign")
		}

		log.Warn().Msgf("Selling %f %s at %f %s", campaign.Volume, event.Product.From, event.Price, event.Product.To)

		campaign.State = entity.CampaignStateBuy

		if err := e.db.Save(campaign); err != nil {
			log.Error().Err(err).Msg("Save Campaign")
		}
	default:
		log.Error().Msgf("campaign sell limit unit (%s) invalid", campaign.SellLimitUnit)
	}
}
*/

func (e *Engine) trade(provider string, event *exchanges.TickerEvent, ts *timeseries.Timeseries) {
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
			// @TODO: use campaign.BuyAlgorithm for get algo
			algo, err := e.algorithms.Get("trend")
			if err != nil {
				log.Error().Err(err).Msg("Get algo failed")

				continue
			}

			algo.Buy(event, campaign, ts)
			// e.tradeBuying(provider, event, campaign, ts)
		} else {
			// @TODO: use campaign.SellAlgorithm for get algo
			algo, err := e.algorithms.Get("trend")
			if err != nil {
				log.Error().Err(err).Msg("Get algo failed")

				continue
			}

			algo.Sell(event, campaign, ts)

			// e.tradeSelling(provider, event, campaign, ts)
		}
	}

	msg := log.Info().
		Str("side", string(event.Side)).
		Int("campaigns", len(campaigns)).
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

					key := fmt.Sprintf("%s-%s", evt.Provider, event.Product.String())

					ts, ok := e.timeseries[key]
					if ok == false {
						log.Error().Msgf("Cannot get timeserie for %s", key)

						continue
					}

					ts.Add(event.Time.Unix(), event.Price)

					e.trade(evt.Provider, event, ts)
				}
			}()

		case evt := <-e.productsCh:
			productsList := []string{}

			for _, product := range evt.Products {
				productsList = append(productsList, product.String())

				key := fmt.Sprintf("%s-%s", evt.Provider, product.String())

				if _, ok := e.timeseries[key]; ok == false {
					e.timeseries[key] = timeseries.New(5000)
				}
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

// GetTimeserie from key {provider}-{product}
func (e *Engine) GetTimeserie(key string) (*timeseries.Timeseries, error) {
	if ts, ok := e.timeseries[key]; ok {
		return ts, nil
	}

	return nil, fmt.Errorf("timeserie %s not found", key)
}
