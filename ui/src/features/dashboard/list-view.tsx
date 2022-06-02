import { FC } from 'react';
import { List } from '@mui/material';
import { useSelector } from 'react-redux';

import { selectCronJobIds } from './state';
import { ListViewItem } from './list-view-item';

export const ListView: FC = () => {
  const cronJobIds = useSelector(selectCronJobIds);

  return (
    <List dense disablePadding>
      {cronJobIds.map((cronJobId) => (
        <ListViewItem key={cronJobId} cronJobId={cronJobId} />
      ))}
    </List>
  );
};
