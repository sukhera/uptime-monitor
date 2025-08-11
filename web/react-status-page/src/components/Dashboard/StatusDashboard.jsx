import { useState } from 'react';
import { useApi } from '../../hooks/useApi';
import { usePolling } from '../../hooks/usePolling';
import ServiceCard from './ServiceCard';
import LoadingSpinner from '../common/LoadingSpinner';
import DesktopSidebar from '../Desktop/DesktopSidebar';

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
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
        <div className="desktop-container desktop-section">
          <div className="glass-card rounded-2xl p-8 max-w-2xl mx-auto text-center">
            <div className="text-6xl mb-4">⚠️</div>
            <h3 className="text-2xl font-semibold text-white mb-4">
              Error Loading Status
            </h3>
            <p className="text-white/70 mb-6 text-lg">
              {error.message || 'Failed to load system status'}
            </p>
            <button
              onClick={() => refetch()}
              className="px-8 py-3 bg-red-500/20 hover:bg-red-500/30 border border-red-500/30 text-red-200 rounded-full transition-all duration-200 desktop-hover-lift"
            >
              🔄 Retry Connection
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      {/* Desktop grid layout */}
      <div className="desktop-container desktop-section">
        <div className="grid grid-cols-12 gap-8">
          
          {/* Desktop sidebar */}
          <div className="col-span-3">
            <DesktopSidebar services={services || []} />
          </div>
          
          {/* Main content area */}
          <div className="col-span-9">
            
            {/* Enhanced header */}
            <div className="mb-12 animate-desktop-fade-in">
              <h1 className="text-4xl lg:text-5xl font-bold text-white mb-4">
                System Status
              </h1>
              <div className="flex items-center space-x-6">
                <p className="text-xl text-white/80">
                  Real-time monitoring of all services
                </p>
                <div className="flex items-center space-x-2 glass-secondary px-4 py-2 rounded-full">
                  <span className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></span>
                  <span className="text-white/70 text-sm">
                    Last updated: {lastUpdated.toLocaleTimeString()}
                  </span>
                </div>
              </div>
            </div>
            
            {/* Enhanced services grid */}
            {!services || services.length === 0 ? (
              <div className="glass-card rounded-2xl p-8 text-center">
                <div className="text-6xl mb-4">📡</div>
                <h3 className="text-2xl font-semibold text-white mb-4">
                  No Services Available
                </h3>
                <p className="text-white/70 text-lg">
                  No services are currently being monitored. Configure your first service to get started.
                </p>
                <button className="mt-6 px-6 py-3 glass-secondary hover:bg-white/10 border border-white/20 text-white rounded-full transition-all duration-200 desktop-hover-lift">
                  ⚙️ Configure Services
                </button>
              </div>
            ) : (
              <div className="grid gap-6 lg:grid-cols-2 xl:grid-cols-3">
                {services.map((service, index) => (
                  <ServiceCard key={service.name || index} service={service} />
                ))}
              </div>
            )}
            
          </div>
        </div>
      </div>
    </div>
  );
};

export default StatusDashboard;