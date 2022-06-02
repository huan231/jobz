import { FC, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Divider, Typography } from '@mui/material';

import { initDashboard, selectInitStatus } from './state';
import { ListView } from './list-view';

export const DashboardPage: FC = () => {
  const dispatch = useDispatch();

  const initStatus = useSelector(selectInitStatus);

  useEffect(() => {
    if (initStatus === 'idle') {
      dispatch(initDashboard());
    }
  }, [initStatus]);

  return (
    <>
      <Typography variant="h5" sx={{ px: 2, py: 2 }}>
        🅹🅾🅱🆉
      </Typography>
      <Divider />
      <ListView />
    </>
  );
};
