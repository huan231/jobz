import { createSelector } from '@reduxjs/toolkit';

import { DashboardState } from './reducer';

const selectDashboardState = (state: { dashboard: DashboardState }) => state.dashboard;

export const selectInitStatus = createSelector([selectDashboardState], (state) => state.initStatus);
export const selectCronJobs = createSelector([selectDashboardState], (state) => state.cronJobs);
export const selectCronJobIds = createSelector([selectDashboardState], (state) => state.cronJobIds);
