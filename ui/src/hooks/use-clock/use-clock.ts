import { useEffect, useState } from 'react';

const ONE_SECOND_IN_MILLISECONDS = 1_000;

export const useClock = () => {
  const [date, setDate] = useState(() => new Date());

  useEffect(() => {
    const interval = setInterval(() => {
      setDate(new Date());
    }, ONE_SECOND_IN_MILLISECONDS);

    return () => {
      clearInterval(interval);
    };
  }, []);

  return date;
};
