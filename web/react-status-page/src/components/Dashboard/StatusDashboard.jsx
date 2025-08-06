import { useState } from 'react';
import { useApi } from '../../hooks/useApi';
import { usePolling } from '../../hooks/usePolling';
import ServiceCard from './ServiceCard';
import LoadingSpinner from '../common/LoadingSpinner';

const StatusDashboard = () => {
  const { data: services, loading, error, refetch } = useApi('/api/v1/status');
  const [lastUpdated, setLastUpdated] = useState(new Date());

  // Auto-refresh every 60 seconds
  usePolling(() => {
    refetch();
    setLastUpdated(new Date());
  }, 60000);

  if (loading) return <LoadingSpinner />;
  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <h3 className="text-lg font-medium text-red-800 dark:text-red-200">
            Error Loading Status
          </h3>
          <p className="text-red-600 dark:text-red-300 mt-1">
            {error.message || 'Failed to load system status'}
          </p>
          <button
            onClick={() => refetch()}
            className="mt-3 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
          System Status
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          Last updated: {lastUpdated.toLocaleTimeString()}
        </p>
      </div>

      {!services || services.length === 0 ? (
        <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
          <h3 className="text-lg font-medium text-yellow-800 dark:text-yellow-200">
            No Services Available
          </h3>
          <p className="text-yellow-600 dark:text-yellow-300 mt-1">
            No services are currently being monitored.
          </p>
        </div>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {services.map((service, index) => (
            <ServiceCard key={service.name || index} service={service} />
          ))}
        </div>
      )}
    </div>
  );
};

export default StatusDashboard;