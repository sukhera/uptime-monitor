const StatusIndicator = ({ status }) => {
  const getStatusConfig = (status) => {
    switch (status?.toLowerCase()) {
      case 'operational':
        return {
          text: 'Operational',
          classes: 'bg-operational-50 text-operational-600 dark:bg-operational-900 dark:text-operational-300'
        };
      case 'degraded':
        return {
          text: 'Degraded',
          classes: 'bg-degraded-50 text-degraded-600 dark:bg-degraded-900 dark:text-degraded-300'
        };
      case 'down':
        return {
          text: 'Down',
          classes: 'bg-down-50 text-down-600 dark:bg-down-900 dark:text-down-300'
        };
      default:
        return {
          text: 'Unknown',
          classes: 'bg-gray-50 text-gray-600 dark:bg-gray-800 dark:text-gray-400'
        };
    }
  };

  const config = getStatusConfig(status);

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.classes}`}>
      {config.text}
    </span>
  );
};

export default StatusIndicator;