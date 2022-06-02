import { KyInstance } from 'ky/distribution/types/ky';

import { CronJobDTO, JobDTO } from '../types';

export interface ApiService {
  getCronJobs: () => Promise<CronJobDTO[]>;
  getJobs: () => Promise<JobDTO[]>;
}

export const makeApiService = (client: KyInstance): ApiService => {
  return {
    getCronJobs: () => client.get('cronjobs').json<CronJobDTO[]>(),
    getJobs: () => client.get('jobs').json<JobDTO[]>(),
  };
};
