import StatusIndicator from './StatusIndicator';

const ServiceCard = ({ service }) => {
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
          <span className="font-medium">{service.latency_ms}ms</span>
        </div>
        <div className="flex justify-between">
          <span>Last Check:</span>
          <span className="font-medium">
            {new Date(service.updated_at).toLocaleTimeString()}
          </span>
        </div>
      </div>
    </div>
  );
};

export default ServiceCard;