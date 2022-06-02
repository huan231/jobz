import { CronJobDTO, RunningJobDTO, SucceededJobDTO, FailedJobDTO } from '../types';

interface BaseEvent<T extends string, P> {
  type: T;
  payload: P;
}

type CronJobAddEvent = BaseEvent<'cronjobadd', CronJobDTO>;
type CronJobDeleteEvent = BaseEvent<'cronjobdelete', Pick<CronJobDTO, 'id'>>;
type CronJobUpdateEvent = BaseEvent<'cronjobdupdate', CronJobDTO>;
type JobAddEvent = BaseEvent<'jobadd', RunningJobDTO>;
type JobCompleteEvent = BaseEvent<'jobcomplete', SucceededJobDTO | FailedJobDTO>;

type OpenHandler = () => void;
type ErrorHandler = () => void;
type MessageHandler = (
  event: CronJobAddEvent | CronJobDeleteEvent | CronJobUpdateEvent | JobAddEvent | JobCompleteEvent,
) => void;

export interface EventsService {
  connect: (handlers: { open: OpenHandler; error: ErrorHandler; message: MessageHandler }) => () => void;
}

export const makeEventsService = (baseUrl: string): EventsService => {
  return {
    connect: (handlers) => {
      const eventSource = new EventSource(new URL('/events', baseUrl));

      eventSource.onopen = () => {
        handlers.open();
      };

      eventSource.onerror = () => {
        handlers.error();
      };

      eventSource.onmessage = (event) => {
        handlers.message(JSON.parse(event.data));
      };

      return () => {
        eventSource.onopen = null;
        eventSource.onerror = null;
        eventSource.onmessage = null;

        eventSource.close();
      };
    },
  };
};
