import StatusIndicator from './StatusIndicator';

const ServiceCard = ({ service }) => {
  // Handle missing or invalid service data
  if (!service || !service.name) {
    return (
      <div className="bg-gray-100 dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div className="text-gray-500 dark:text-gray-400">
          Invalid service data
        </div>
      </div>
    );
  }

  const formatLatency = (latency) => {
    if (latency === undefined || latency === null) return 'N/A';
    return `${latency}ms`;
  };

  const formatTimestamp = (timestamp) => {
    if (!timestamp) return 'N/A';
    try {
      return new Date(timestamp).toLocaleTimeString();
    } catch {
      return 'Invalid date';
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 transition-all duration-200 hover:shadow-lg">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          {service.name}
        </h3>
        <StatusIndicator status={service.status} />
      </div>
      
      <div className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
        <div className="flex justify-between">
          <span>Response Time:</span>
          <span className="font-medium">{formatLatency(service.latency_ms)}</span>
        </div>
        <div className="flex justify-between">
          <span>Last Check:</span>
          <span className="font-medium">
            {formatTimestamp(service.updated_at)}
          </span>
        </div>
        {service.error && (
          <div className="mt-2 p-2 bg-red-50 dark:bg-red-900/20 rounded text-xs text-red-600 dark:text-red-300">
            Error: {service.error}
          </div>
        )}
      </div>
    </div>
  );
};

export default ServiceCard;