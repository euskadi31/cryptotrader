
/*
{
    "product":"BTC-EUR",
    "price":13470.1,
    "side":"sell",
    "time":"2017-12-22T00:09:24.015Z",
    "size":0.00371192
}
*/
export class TickerEvent {
    constructor(
        public product?: string,
        public price?: number,
        public side?: string,
        public time?: string,
        public size?: number,
    ) {}
}
