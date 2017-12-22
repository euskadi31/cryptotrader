
export class Campaign {
    constructor(
        public id?: number,
        public provider?: string,
        public product_id?: string,
        public volume?: number,
        public buy_limit?: number,
        public sell_limit?: number,
        public sell_limit_unit?: string,
        public created_at?: string,
        public updated_at?: string,
        public state?: string,
    ) {}
}
