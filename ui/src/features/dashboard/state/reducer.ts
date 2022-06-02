import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { differenceInMilliseconds, parseISO } from 'date-fns';

import { CronJobDTO, FailedJobDTO, JobDTO, RunningJobDTO, SucceededJobDTO } from '../types';

export type CompletedJob = (FailedJobDTO | SucceededJobDTO) & { duration: number };

export interface CronJob extends CronJobDTO {
  status: 'pending' | JobDTO['status'];
  jobs: (RunningJobDTO | CompletedJob)[];
}

const makeJobDuration = (job: FailedJobDTO | SucceededJobDTO) => {
  return differenceInMilliseconds(parseISO(job.completedAt), parseISO(job.createdAt));
};

export interface DashboardState {
  initStatus: 'idle' | 'fetching' | 'succeeded' | 'failed';
  streamStatus: 'idle' | 'connecting' | 'connected' | 'disconnected';
  cronJobIds: CronJob['id'][];
  cronJobs: Record<CronJob['id'], CronJob>;
}

const initialState: DashboardState = {
  initStatus: 'idle',
  streamStatus: 'idle',
  cronJobIds: [],
  cronJobs: {},
};

const dashboardSlice = createSlice({
  name: 'dashboard',
  initialState,
  reducers: {
    initDashboard: (state) => {
      state.initStatus = 'fetching';
    },
    initDashboardSuccess: (state, action: PayloadAction<{ cronJobs: CronJobDTO[]; jobs: JobDTO[] }>) => {
      state.initStatus = 'succeeded';
      state.cronJobIds = action.payload.cronJobs.map((cronJob) => cronJob.id);
      state.cronJobs = action.payload.cronJobs.reduce<DashboardState['cronJobs']>((cronJobs, cronJob) => {
        const jobs = action.payload.jobs
          .filter((job) => job.cronJobId === cronJob.id)
          .map((job) => (job.status === 'running' ? job : { ...job, duration: makeJobDuration(job) }));

        cronJobs[cronJob.id] = { ...cronJob, status: jobs[jobs.length - 1]?.status ?? 'pending', jobs };

        return cronJobs;
      }, {});
    },
    initDashboardFailure: (state) => {
      state.initStatus = 'failed';
    },
    cronJobAdded: (state, action: PayloadAction<CronJobDTO>) => {
      state.cronJobIds.push(action.payload.id);

      state.cronJobs[action.payload.id] = { ...action.payload, status: 'pending', jobs: [] };
    },
    cronJobDeleted: (state, action: PayloadAction<Pick<CronJobDTO, 'id'>>) => {
      delete state.cronJobs[action.payload.id];

      const index = state.cronJobIds.findIndex((cronJobId) => cronJobId === action.payload.id);

      if (index !== -1) {
        state.cronJobIds.splice(index, 1);
      }
    },
    cronJobUpdated: (state, action: PayloadAction<CronJobDTO>) => {
      if (typeof state.cronJobs[action.payload.id] !== 'undefined') {
        state.cronJobs[action.payload.id] = { ...state.cronJobs[action.payload.id], ...action.payload };
      }
    },
    jobAdded: (state, action: PayloadAction<RunningJobDTO>) => {
      const { cronJobId } = action.payload;

      if (typeof state.cronJobs[cronJobId] !== 'undefined') {
        state.cronJobs[cronJobId].status = action.payload.status;
        state.cronJobs[cronJobId].jobs.push(action.payload);
      }
    },
    jobCompleted: (state, action: PayloadAction<SucceededJobDTO | FailedJobDTO>) => {
      const { cronJobId } = action.payload;

      if (typeof state.cronJobs[cronJobId] === 'undefined') {
        return;
      }

      const index = state.cronJobs[cronJobId].jobs.findIndex((job) => job.id === action.payload.id);

      if (index !== -1) {
        state.cronJobs[cronJobId].status = action.payload.status;
        state.cronJobs[cronJobId].jobs[index] = { ...action.payload, duration: makeJobDuration(action.payload) };
      }
    },
    streamConnecting: (state) => {
      state.streamStatus = 'connecting';
    },
    streamConnected: (state) => {
      state.streamStatus = 'connected';
    },
    streamDisconnected: (state) => {
      state.streamStatus = 'disconnected';
    },
  },
});

export const dashboardReducer = dashboardSlice.reducer;
export const {
  initDashboard,
  initDashboardSuccess,
  initDashboardFailure,
  cronJobAdded,
  cronJobUpdated,
  cronJobDeleted,
  jobAdded,
  jobCompleted,
  streamConnecting,
  streamConnected,
  streamDisconnected,
} = dashboardSlice.actions;
