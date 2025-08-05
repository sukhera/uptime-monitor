import { useState, useEffect, useCallback } from 'react';
import axios from 'axios';

// Create axios instance with default configuration
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor for logging
apiClient.interceptors.request.use(
  (config) => {
    console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('[API] Request error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('[API] Response error:', error);
    
    // Handle specific error cases
    if (error.response?.status === 429) {
      console.warn('[API] Rate limit exceeded');
    } else if (error.response?.status >= 500) {
      console.error('[API] Server error:', error.response.status);
    }
    
    return Promise.reject(error);
  }
);

/**
 * Custom hook for API calls with retry logic and error handling
 * @param {string} url - API endpoint
 * @param {Object} options - Configuration options
 * @param {number} options.retries - Number of retry attempts (default: 3)
 * @param {number} options.retryDelay - Delay between retries in ms (default: 1000)
 * @param {boolean} options.autoRetry - Whether to automatically retry on failure (default: true)
 * @param {boolean} options.polling - Whether to poll the endpoint (default: false)
 * @param {number} options.pollingInterval - Polling interval in ms (default: 30000)
 */
export const useApi = (url, options = {}) => {
  const {
    retries = 3,
    retryDelay = 1000,
    autoRetry = true,
    polling = false,
    pollingInterval = 30000,
    ...axiosOptions
  } = options;

  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [retryCount, setRetryCount] = useState(0);

  const executeRequest = useCallback(async (isRetry = false) => {
    if (!isRetry) {
      setLoading(true);
      setError(null);
    }

    try {
      const response = await apiClient.get(url, axiosOptions);
      
      // Handle successful response
      setData(response.data);
      setError(null);
      setRetryCount(0);
      
      return response.data;
    } catch (err) {
      const shouldRetry = autoRetry && retryCount < retries && err.response?.status >= 500;
      
      if (shouldRetry) {
        console.warn(`[API] Retry ${retryCount + 1}/${retries} for ${url}`);
        setRetryCount(prev => prev + 1);
        
        // Wait before retrying
        await new Promise(resolve => setTimeout(resolve, retryDelay));
        return executeRequest(true);
      } else {
        setError(err);
        setLoading(false);
        throw err;
      }
    } finally {
      if (!isRetry) {
        setLoading(false);
      }
    }
  }, [url, retries, retryDelay, autoRetry, retryCount, axiosOptions]);

  // Polling effect
  useEffect(() => {
    if (!polling) return;

    const interval = setInterval(() => {
      executeRequest();
    }, pollingInterval);

    return () => clearInterval(interval);
  }, [polling, pollingInterval, executeRequest]);

  // Manual retry function
  const retry = useCallback(() => {
    setRetryCount(0);
    return executeRequest();
  }, [executeRequest]);

  // Refresh function
  const refresh = useCallback(() => {
    return executeRequest();
  }, [executeRequest]);

  return {
    data,
    loading,
    error,
    retryCount,
    retry,
    refresh,
    executeRequest,
  };
};

/**
 * Hook for fetching status data with automatic polling
 */
export const useStatusApi = () => {
  return useApi('/status', {
    polling: true,
    pollingInterval: 30000, // 30 seconds
    retries: 3,
    retryDelay: 2000,
  });
};

/**
 * Hook for health check endpoint
 */
export const useHealthApi = () => {
  return useApi('/health', {
    polling: false,
    retries: 1,
    retryDelay: 1000,
  });
};