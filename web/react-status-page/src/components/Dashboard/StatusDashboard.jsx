import { useState } from 'react';
import { useApi } from '../../hooks/useApi';
import { usePolling } from '../../hooks/usePolling';
import ServiceCard from './ServiceCard';
import LoadingSpinner from '../common/LoadingSpinner';

const StatusDashboard = () => {
  const { data: services, loading, error, refetch } = useApi('/api/status');
  const [lastUpdated, setLastUpdated] = useState(new Date());

  // Auto-refresh every 60 seconds
  usePolling(() => {
    refetch();
    setLastUpdated(new Date());
  }, 60000);

  if (loading) return <LoadingSpinner />;
  if (error) return <div className="text-red-500">Error: {error}</div>;

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

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {services?.map(service => (
          <ServiceCard key={service.slug} service={service} />
        ))}
      </div>
    </div>
  );
};

export default StatusDashboard;