import { Injectable } from '@angular/core';

import { Campaign } from '../entities/campaign';
import { ApiService } from './api.service';

@Injectable()
export class CampaignService {
    constructor(private apiService: ApiService) { }

    getCampaigns(): Promise<Campaign[]> {
        return this.apiService.get('/v1/campaigns')
            .then(response => {
                return response.json().map(item => {
                    return Object.assign(new Campaign(), item);
                });
            });
    }
}
