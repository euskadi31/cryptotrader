import { Component, OnInit, OnDestroy, AfterContentInit, ViewChild, ElementRef } from '@angular/core';
import { State } from '@clr/angular';
import { Subscription } from 'rxjs/Subscription';
import { Campaign } from '../../entities/campaign';
import { CampaignService } from '../../services/campaign.service';
import { TickerService } from '../../services/ticker.service';

import * as Highcharts from 'highcharts';

@Component({
    selector: 'app-campaign',
    templateUrl: './campaign.component.html',
    styleUrls: ['./campaign.component.less']
})
export class CampaignComponent implements OnInit, AfterContentInit, OnDestroy {
    campaigns: Campaign[];
    total: number;
    loading: boolean;

    chart: Highcharts.ChartObject;

    @ViewChild('chart')
    chartTarget: ElementRef;

    private subscriptions: Subscription[];

    constructor(private campaignService: CampaignService, private tickerService: TickerService) {
        this.loading = true;
        this.subscriptions = [];
    }

    ngOnInit() {
        this.subscriptions.push(
            this.tickerService.ticker('gdax', 'btc-eur').subscribe(ticker => {
                this.chart.series[0].addPoint([(new Date(ticker.time)).getTime(), ticker.price]);
            })
        );
    }

    ngOnDestroy() {
        this.subscriptions.forEach(subscription => subscription.unsubscribe());
        this.chart = null;
    }

    ngAfterContentInit() {
        const options: Highcharts.Options = {
            chart: {
                zoomType: 'x'
            },
            title: {
                text: 'BTC to EUR exchange rate over time'
            },
            subtitle: {
                text: document.ontouchstart === undefined ?
                        'Click and drag in the plot area to zoom in' : 'Pinch the chart to zoom in'
            },
            xAxis: {
                type: 'datetime',
                plotLines: [
                    {
                        value: 1,
                        color: '#00FF00',
                        dashStyle: 'shortdash',
                        width: 1,
                    },
                    {
                        value: 8,
                        color: '#FF00FF',
                        dashStyle: 'shortdash',
                        width: 1,
                    }
                ]
            },
            yAxis: {
                title: {
                    text: 'Price'
                }
            },
            legend: {
                enabled: false
            },
            credits: {
                enabled: false
            },
            plotOptions: {
                area: {
                    fillColor: {
                        linearGradient: {
                            x1: 0,
                            y1: 0,
                            x2: 0,
                            y2: 1
                        },
                        stops: [
                            [0, Highcharts.getOptions().colors[0]],
                            [1, (<any>Highcharts.Color(Highcharts.getOptions().colors[0])).setOpacity(0).get('rgba')]
                        ]
                    },
                    marker: {
                        radius: 2
                    },
                    lineWidth: 1,
                    states: {
                        hover: {
                            lineWidth: 1
                        }
                    },
                    threshold: null
                }
            },

            series: [
                {
                    type: 'area',
                    name: 'USD to EUR',
                    data: [
                    ]
                },
                {
                    type: 'scatter',
                    name: 'Campaign 1',
                    data: [
                        /*{
                            x: 1,
                            y: 15002.540000,
                            name: 'Sell',
                            color: '#00FF00'
                        },
                        {
                            x: 8,
                            y: 15020.880000,
                            name: 'Buy',
                            color: '#FF00FF'
                        }*/
                    ]
                }
            ]
        };

        this.chart = Highcharts.chart(this.chartTarget.nativeElement, options);
    }

    fetchCampaign(filters: {[prop: string]: any}) {
        this.campaignService.getCampaigns().then(campaigns => {
            this.total = campaigns.length;
            this.campaigns = campaigns;

            this.loading = false;
        });
    }

    refresh(state: State) {
        this.loading = true;
        // We convert the filters from an array to a map,
        // because that's what our backend-calling service is expecting

        const filters: {[prop: string]: any[]} = {};

        if (state.filters) {
            for (const filter of state.filters) {
                const {property, value} = <{property: string, value: string}>filter;
                filters[property] = [value];
            }
        }

        this.fetchCampaign(filters);
    }
}
