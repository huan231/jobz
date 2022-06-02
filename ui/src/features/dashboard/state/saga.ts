import { all, call, put, take, takeLatest, takeLeading } from 'redux-saga/effects';
import { END, eventChannel } from 'redux-saga';

import {
  cronJobAdded,
  cronJobDeleted,
  cronJobUpdated,
  initDashboard,
  initDashboardFailure,
  initDashboardSuccess,
  jobAdded,
  jobCompleted,
  streamConnected,
  streamConnecting,
  streamDisconnected,
} from './reducer';
import { ApiService, EventsService } from '../services';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Unwrap<T extends (...args: any[]) => any> = Awaited<ReturnType<T>>;

export const makeInitDashboardSaga = (api: ApiService) => {
  return function* () {
    try {
      const [cronJobs, jobs]: [Unwrap<ApiService['getCronJobs']>, Unwrap<ApiService['getJobs']>] = yield all([
        call(api.getCronJobs),
        call(api.getJobs),
      ]);

      yield put(initDashboardSuccess({ cronJobs, jobs }));
    } catch (err) {
      yield put(initDashboardFailure());
    }
  };
};

type EventAction = ReturnType<
  | typeof cronJobAdded
  | typeof cronJobDeleted
  | typeof cronJobUpdated
  | typeof jobAdded
  | typeof jobCompleted
  | typeof streamConnected
>;

const createEventsChannel = (events: EventsService) => {
  return eventChannel<EventAction>((emit) => {
    const disconnect = events.connect({
      open: () => {
        emit(streamConnected());
      },
      error: () => {
        emit(END);
      },
      message: (event) => {
        switch (event.type) {
          case 'cronjobadd':
            return emit(cronJobAdded(event.payload));
          case 'cronjobdelete':
            return emit(cronJobDeleted(event.payload));
          case 'cronjobdupdate':
            return emit(cronJobUpdated(event.payload));
          case 'jobadd':
            return emit(jobAdded(event.payload));
          case 'jobcomplete':
            return emit(jobCompleted(event.payload));
        }
      },
    });

    return () => {
      disconnect();
    };
  });
};

export const makeInitDashboardSuccessSaga = (events: EventsService) => {
  return function* () {
    try {
      yield put(streamConnecting());

      const channel: ReturnType<typeof createEventsChannel> = yield call(createEventsChannel, events);

      while (true) {
        const action: EventAction = yield take(channel);

        yield put(action);
      }
    } finally {
      yield put(streamDisconnected());
    }
  };
};

export function* dashboardSaga({ api, events }: { api: ApiService; events: EventsService }) {
  yield* [
    takeLatest(initDashboard, makeInitDashboardSaga(api)),
    takeLeading(initDashboardSuccess, makeInitDashboardSuccessSaga(events)),
  ];
}
