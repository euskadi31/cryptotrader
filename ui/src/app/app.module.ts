import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';
import { ClarityModule } from '@clr/angular';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './components/app/app.component';
import { CampaignComponent } from './components/campaign/campaign.component';
import { NotFoundComponent } from './components/not-found/not-found.component';
import { ApiService } from './services/api.service';
import { CampaignService } from './services/campaign.service';
import { TickerService } from './services/ticker.service';

@NgModule({
    declarations: [
        AppComponent,
        CampaignComponent,
        NotFoundComponent
    ],
    imports: [
        BrowserModule,
        HttpModule,
        AppRoutingModule,
        ClarityModule,
        BrowserAnimationsModule,
    ],
    providers: [
        ApiService,
        CampaignService,
        TickerService
    ],
    bootstrap: [
        AppComponent
    ]
})
export class AppModule { }
