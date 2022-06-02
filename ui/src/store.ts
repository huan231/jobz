import { configureStore } from '@reduxjs/toolkit';
import createSagaMiddleware from 'redux-saga';
import { all, fork } from 'redux-saga/effects';

import { dashboardReducer, dashboardSaga } from './features';
import { dashboardApi, dashboardEvents } from './services';

const sagaMiddleware = createSagaMiddleware();

export const store = configureStore({
  reducer: { dashboard: dashboardReducer },
  middleware: (getDefaultMiddleware) => [...getDefaultMiddleware({ thunk: false }), sagaMiddleware],
});

sagaMiddleware.run(function* () {
  yield all([fork(dashboardSaga, { api: dashboardApi, events: dashboardEvents })]);
});
