import { format, isAfter } from 'date-fns';

const MaintenanceSchedule = ({ maintenance }) => {
  const upcomingMaintenance = maintenance?.filter(
    item => isAfter(new Date(item.scheduled_start), new Date())
  );

  if (!upcomingMaintenance?.length) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 dark:text-gray-400">
          No scheduled maintenance at this time.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {upcomingMaintenance.map(item => (
        <div key={item.id} className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
          <div className="flex items-start justify-between mb-4">
            <div>
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                {item.title}
              </h3>
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                {format(new Date(item.scheduled_start), 'PPp')} - {format(new Date(item.scheduled_end), 'p')}
              </p>
            </div>
            <span className="px-2 py-1 bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300 rounded-full text-xs font-medium">
              Scheduled
            </span>
          </div>
          
          <p className="text-gray-700 dark:text-gray-300 mb-4">
            {item.description}
          </p>
          
          {item.impact && (
            <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md p-3">
              <p className="text-yellow-800 dark:text-yellow-300 text-sm">
                <strong>Expected Impact:</strong> {item.impact}
              </p>
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default MaintenanceSchedule;