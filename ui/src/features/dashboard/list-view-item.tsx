import { FC, PropsWithChildren, ReactElement, useMemo } from 'react';
import cronstrue from 'cronstrue';
import { parseCronExpression } from 'cron-schedule';
import { Box, ListItem, ListItemText, Typography } from '@mui/material';
import {
  AccessTime,
  AccessTimeOutlined,
  CheckCircleOutlineOutlined,
  ErrorOutlineOutlined,
  Favorite,
  SvgIconComponent,
  Timelapse,
  Timer,
  Workspaces,
} from '@mui/icons-material';
import { formatDuration, intervalToDuration } from 'date-fns';
import { createSelector } from '@reduxjs/toolkit';
import { useSelector } from 'react-redux';

import { CompletedJob, CronJob, selectCronJobs } from './state';
import { useClock } from '../../hooks';

const STATUS_ICONS: Record<CronJob['status'], ReactElement> = {
  pending: <AccessTimeOutlined color="disabled" fontSize="small" />,
  running: <AccessTimeOutlined color="info" fontSize="small" />,
  succeeded: <CheckCircleOutlineOutlined color="success" fontSize="small" />,
  failed: <ErrorOutlineOutlined color="error" fontSize="small" />,
};

interface ListViewItemProps {
  cronJobId: CronJob['id'];
}

const makeSelectors = (cronJobId: CronJob['id']) => {
  const selectCronJob = createSelector([selectCronJobs], (cronJobs) => cronJobs[cronJobId]);

  return {
    selectNamespace: createSelector([selectCronJob], (cronJob) => cronJob.namespace),
    selectName: createSelector([selectCronJob], (cronJob) => cronJob.name),
    selectSchedule: createSelector([selectCronJob], (cronJob) => cronJob.schedule),
    selectStatus: createSelector([selectCronJob], (cronJob) => cronJob.status),
    selectCompletedJobs: createSelector([selectCronJob], (cronJob) =>
      cronJob.jobs.filter((job) => job.status !== 'running'),
    ),
    selectSucceededJobs: createSelector([selectCronJob], (cronJob) =>
      cronJob.jobs.filter((job): job is CompletedJob => job.status === 'succeeded'),
    ),
  };
};

interface ListViewItemAdornmentProps {
  icon: SvgIconComponent;
}

const ListViewItemAdornment: FC<PropsWithChildren<ListViewItemAdornmentProps>> = ({ icon: Icon, children }) => {
  return (
    <Box sx={{ display: 'flex', alignItems: 'center', columnGap: 1 }}>
      <Icon fontSize="small" />
      <Typography variant="caption" noWrap sx={{ lineHeight: '1.5' }}>
        {children}
      </Typography>
    </Box>
  );
};

interface ListViewItemSuccessRateProps {
  completedJobsCount: number;
  succeededJobsCount: number;
}

const ListViewItemSuccessRate: FC<ListViewItemSuccessRateProps> = ({ completedJobsCount, succeededJobsCount }) => {
  const text = useMemo(() => {
    if (completedJobsCount === 0) {
      return 'n/a';
    }

    return `${Math.round((succeededJobsCount / completedJobsCount) * 100)}%`;
  }, [completedJobsCount, succeededJobsCount]);

  return <ListViewItemAdornment icon={Favorite}>{text}</ListViewItemAdornment>;
};

interface ListViewItemAvgDurationProps {
  succeededJobs: CompletedJob[];
}

const ListViewItemAvgDuration: FC<ListViewItemAvgDurationProps> = ({ succeededJobs }) => {
  const text = useMemo(() => {
    if (succeededJobs.length === 0) {
      return 'n/a';
    }

    const duration = succeededJobs.reduce((duration, job) => duration + job.duration, 0);

    return formatDuration(intervalToDuration({ start: 0, end: Math.round(duration / succeededJobs.length) }));
  }, [succeededJobs]);

  return <ListViewItemAdornment icon={Timer}>{text}</ListViewItemAdornment>;
};

interface ListViewItemNextExecutionProps {
  schedule: CronJob['schedule'];
}

const ListViewItemNextExecution: FC<ListViewItemNextExecutionProps> = ({ schedule }) => {
  const cron = useMemo(() => parseCronExpression(schedule), [schedule]);
  const date = useClock();

  const text = useMemo(() => {
    const formatted = formatDuration(intervalToDuration({ start: date, end: cron.getNextDate() }));

    if (formatted === '') {
      return 'now';
    }

    return `in ${formatted}`;
  }, [cron, date]);

  return <ListViewItemAdornment icon={Timelapse}>{text}</ListViewItemAdornment>;
};

export const ListViewItem: FC<ListViewItemProps> = ({ cronJobId }) => {
  const { selectNamespace, selectName, selectSchedule, selectStatus, selectCompletedJobs, selectSucceededJobs } =
    useMemo(() => makeSelectors(cronJobId), [cronJobId]);

  const namespace = useSelector(selectNamespace);
  const name = useSelector(selectName);
  const schedule = useSelector(selectSchedule);
  const status = useSelector(selectStatus);
  const completedJobs = useSelector(selectCompletedJobs);
  const succeededJobs = useSelector(selectSucceededJobs);

  const description = useMemo(() => cronstrue.toString(schedule).toLowerCase(), [schedule]);

  return (
    <ListItem divider>
      <ListItemText
        disableTypography
        primary={
          <Box sx={{ display: 'flex', alignItems: 'center', columnGap: 1 }}>
            {STATUS_ICONS[status]}
            <Typography variant="h6">{name}</Typography>
          </Box>
        }
        secondary={
          <Box
            sx={{
              display: 'flex',
              mt: 1,
              columnGap: { sm: 1, md: 2 },
              color: 'text.secondary',
              flexWrap: 'wrap',
              flexDirection: { xs: 'column', sm: 'row' },
              rowGap: 0.5,
            }}
          >
            <ListViewItemAdornment icon={Workspaces}>{namespace}</ListViewItemAdornment>
            <ListViewItemAdornment icon={AccessTime}>{description}</ListViewItemAdornment>
            <ListViewItemSuccessRate
              completedJobsCount={completedJobs.length}
              succeededJobsCount={succeededJobs.length}
            />
            <ListViewItemAvgDuration succeededJobs={succeededJobs} />
            <ListViewItemNextExecution schedule={schedule} />
          </Box>
        }
      />
    </ListItem>
  );
};
