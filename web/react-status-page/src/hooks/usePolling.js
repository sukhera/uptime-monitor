import { useEffect, useRef } from 'react';

export const usePolling = (callback, interval = 60000) => {
  const intervalRef = useRef();

  useEffect(() => {
    if (callback && interval > 0) {
      intervalRef.current = setInterval(callback, interval);
      return () => clearInterval(intervalRef.current);
    }
  }, [callback, interval]);

  const startPolling = () => {
    if (!intervalRef.current) {
      intervalRef.current = setInterval(callback, interval);
    }
  };

  const stopPolling = () => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  };

  return { startPolling, stopPolling };
};