import ky from 'ky';

import { makeApiService, makeEventsService } from './features/dashboard/services';
import { config } from './config';

const client = ky.create({ prefixUrl: config.apiBaseUrl });

export const dashboardApi = makeApiService(client);
export const dashboardEvents = makeEventsService(config.apiBaseUrl);
