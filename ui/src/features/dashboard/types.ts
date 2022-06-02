export interface CronJobDTO {
  id: string;
  namespace: string;
  name: string;
  schedule: string;
}

interface BaseJobDTO {
  id: string;
  cronJobId: string;
  createdAt: string;
  updatedAt: string;
}

export interface RunningJobDTO extends BaseJobDTO {
  status: 'running';
  completedAt: null;
}

export interface SucceededJobDTO extends BaseJobDTO {
  status: 'succeeded';
  completedAt: string;
}

export interface FailedJobDTO extends BaseJobDTO {
  status: 'failed';
  completedAt: string;
}

export type JobDTO = RunningJobDTO | SucceededJobDTO | FailedJobDTO;
