import { formatDistanceToNow } from 'date-fns';

const IncidentCard = ({ incident }) => {
  const getSeverityColor = (severity) => {
    switch (severity?.toLowerCase()) {
      case 'critical':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300';
      case 'major':
        return 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300';
      case 'minor':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300';
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            {incident.title}
          </h3>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            {formatDistanceToNow(new Date(incident.created_at))} ago
          </p>
        </div>
        <span className={`px-2 py-1 rounded-full text-xs font-medium ${getSeverityColor(incident.severity)}`}>
          {incident.severity}
        </span>
      </div>
      
      <p className="text-gray-700 dark:text-gray-300 mb-4">
        {incident.description}
      </p>
      
      {incident.affected_services && (
        <div className="flex flex-wrap gap-2">
          {incident.affected_services.map(service => (
            <span 
              key={service}
              className="px-2 py-1 bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300 rounded text-xs"
            >
              {service}
            </span>
          ))}
        </div>
      )}
    </div>
  );
};

export default IncidentCard;