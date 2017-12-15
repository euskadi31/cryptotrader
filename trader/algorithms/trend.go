// Copyright 2017 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algorithms

import (
	"encoding/json"

	"github.com/euskadi31/cryptotrader/database/entity"
	"github.com/euskadi31/cryptotrader/exchanges"
	"github.com/euskadi31/cryptotrader/services"
	"github.com/euskadi31/cryptotrader/timeseries"
	"github.com/rs/zerolog/log"
)

const (
	TrendSellingLongTrendSize  = "trend.selling.long_trend_size"
	TrendSellingShortTrendSize = "trend.selling.short_trend_size"
)

// Trend struct
type Trend struct {
	campaignService services.CampaignServiceSave
	orderService    services.OrderServiceSave
}

// NewTrend algorithms
func NewTrend(campaignService services.CampaignServiceSave, orderService services.OrderServiceSave) *Trend {
	return &Trend{
		campaignService: campaignService,
		orderService:    orderService,
	}
}

// Name implements Algorithm interface
func (a Trend) Name() string {
	return "trend"
}

// Options implements Algorithm interface
func (a Trend) Options() Options {
	return Options{
		TrendSellingLongTrendSize:  150,
		TrendSellingShortTrendSize: 10,
	}
}

// MarshalJSON implements json.Marshaler.
func (a Trend) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Options())
}

// Buy implements Algorithm interface
func (a *Trend) Buy(event *exchanges.TickerEvent, campaign *entity.Campaign, ts *timeseries.Timeseries) {
	if event.Price >= campaign.BuyLimit {
		return
	}

	if event.Price >= campaign.BuyLimit {
		return
	}

	options := a.Options()
	options.Merge(campaign.BuyAlgorithmOptions)

	campaign.State = entity.CampaignStateBuying

	if err := a.campaignService.Save(campaign); err != nil {
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

	if err := a.orderService.Save(order); err != nil {
		log.Error().Err(err).Msg("save order failed")

		return
	}

	campaign.State = entity.CampaignStateSell

	campaign.BuyOrder = order
	// campaign.AddOrder(order)

	if err := a.campaignService.Save(campaign); err != nil {
		log.Error().Err(err).Msg("Save Campaign")
	}
	// end emulate

}

// Sell implements Algorithm interface
func (a *Trend) Sell(event *exchanges.TickerEvent, campaign *entity.Campaign, ts *timeseries.Timeseries) {
	switch campaign.SellLimitUnit {
	case "percent":
		log.Debug().Msgf("Current Price: %v", event.Price)
		log.Debug().
			Float64("buy_price", campaign.BuyOrder.GetBuyingMarketPrice()).
			Float64("current_price", campaign.BuyOrder.GetCurrentPrice(event.Price)).
			Msgf("Margin in %%: %v", campaign.BuyOrder.GetMarginInPercent(event.Price))

		if campaign.BuyOrder.GetMarginInPercent(event.Price) < campaign.SellLimit {
			return
		}
	case "currency":
		log.Debug().Msgf("Current Price: %v", event.Price)
		log.Debug().Msgf("Margin in €: %v", campaign.BuyOrder.GetMarginInCurrency(event.Price))

		if campaign.BuyOrder.GetMarginInCurrency(event.Price) < campaign.SellLimit {
			return
		}
	default:
		log.Error().Msgf("campaign sell limit unit (%s) invalid", campaign.SellLimitUnit)

		return
	}

	options := a.Options()
	options.Merge(campaign.SellAlgorithmOptions)

	longTrendSize := options.GetInt(TrendSellingLongTrendSize)
	shortTrendSize := options.GetInt(TrendSellingShortTrendSize)

	historyMaxSize := longTrendSize

	if historyMaxSize < shortTrendSize {
		historyMaxSize = shortTrendSize
	}

	if ts.Size() < historyMaxSize {
		log.Debug().Msg("there are not enough elements in the time series")

		return
	}

	longTrend, err := ts.GetTrending(longTrendSize)
	if err != nil {
		log.Error().Err(err).Msg("GetTrending failed")

		return
	}

	shortTrend, err := ts.GetTrending(shortTrendSize)
	if err != nil {
		log.Error().Err(err).Msg("GetTrending failed")

		return
	}

	if longTrend != timeseries.TrendTypeIncreasing && shortTrend != timeseries.TrendTypeDecreasing {
		log.Debug().Msg("Not match trend model")

		return
	}

	campaign.State = entity.CampaignStateSelling

	if err := a.campaignService.Save(campaign); err != nil {
		log.Error().Err(err).Msg("Save Campaign")
	}

	log.Warn().Msgf("Selling %f %s at %f %s", campaign.Volume, event.Product.From, event.Price, event.Product.To)

	campaign.State = entity.CampaignStateBuy

	if err := a.campaignService.Save(campaign); err != nil {
		log.Error().Err(err).Msg("Save Campaign")
	}
}
