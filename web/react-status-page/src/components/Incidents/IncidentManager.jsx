import { useState } from 'react';
import { useApi } from '../../hooks/useApi';
import IncidentCard from './IncidentCard';
import MaintenanceSchedule from './MaintenanceSchedule';

const IncidentManager = () => {
  const [activeTab, setActiveTab] = useState('incidents');
  const { data: incidents } = useApi('/api/incidents');
  const { data: maintenance } = useApi('/api/maintenance');

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
          Incidents & Maintenance
        </h2>
        
        <div className="border-b border-gray-200 dark:border-gray-700">
          <nav className="-mb-px flex space-x-8">
            <button
              onClick={() => setActiveTab('incidents')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'incidents'
                  ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400'
              }`}
            >
              Recent Incidents
            </button>
            <button
              onClick={() => setActiveTab('maintenance')}
              className={`py-2 px-1 border-b-2 font-medium text-sm ${
                activeTab === 'maintenance'
                  ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400'
              }`}
            >
              Maintenance Schedule
            </button>
          </nav>
        </div>
      </div>

      <div className="mt-6">
        {activeTab === 'incidents' && (
          <div className="space-y-4">
            {incidents?.map(incident => (
              <IncidentCard key={incident.id} incident={incident} />
            ))}
          </div>
        )}
        
        {activeTab === 'maintenance' && (
          <MaintenanceSchedule maintenance={maintenance} />
        )}
      </div>
    </div>
  );
};

export default IncidentManager;