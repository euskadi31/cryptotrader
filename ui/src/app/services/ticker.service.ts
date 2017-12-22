import { NgZone, Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';

import { TickerEvent } from '../entities/ticker';
import { environment } from '../../environments/environment';


@Injectable()
export class TickerService {
    private base: string;

    private eventSource: any = window['EventSource'];

    constructor(private ngZone: NgZone) {
        this.base = environment.api + '/api';
    }

    ticker(provider: string, product: string): Observable<TickerEvent> {
        return new Observable<TickerEvent>(obs => {
            product = product.toLowerCase();

            const eventSource = new this.eventSource(`${this.base}/v1/timeseries/${provider}/${product}/events`);

            eventSource.onmessage = event => {
                const data = JSON.parse(event.data);

                this.ngZone.run(() => obs.next(Object.assign(new TickerEvent(), data)));
            };

            return () => eventSource.close();
        });
    }
}
