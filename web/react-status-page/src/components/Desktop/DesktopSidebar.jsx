const DesktopSidebar = ({ services }) => {
  const totalServices = services?.length || 0;
  const operationalCount = services?.filter(s => s.status === 'operational').length || 0;
  const degradedCount = services?.filter(s => s.status === 'degraded').length || 0;
  const downCount = services?.filter(s => s.status === 'down').length || 0;
  
  // Calculate overall system status
  const getOverallStatus = () => {
    if (downCount > 0) return { status: 'down', text: 'Issues Detected', color: 'text-red-400' };
    if (degradedCount > 0) return { status: 'degraded', text: 'Degraded Performance', color: 'text-yellow-400' };
    return { status: 'operational', text: 'All Systems Operational', color: 'text-green-400' };
  };

  const overallStatus = getOverallStatus();
  
  // Get unique categories if available
  const categories = services
    ? [...new Set(services.map(s => s.category).filter(Boolean))]
    : [];

  const getCategoryStatus = (category) => {
    const categoryServices = services.filter(s => s.category === category);
    const allOperational = categoryServices.every(s => s.status === 'operational');
    const hasDown = categoryServices.some(s => s.status === 'down');
    
    if (hasDown) return { color: 'text-red-400', symbol: '●' };
    if (!allOperational) return { color: 'text-yellow-400', symbol: '●' };
    return { color: 'text-green-400', symbol: '●' };
  };

  return (
    <div className="sticky top-8 space-y-6">
      
      {/* System overview */}
      <div className="glass-card rounded-2xl p-6 animate-desktop-fade-in">
        <h3 className="text-xl font-semibold text-white mb-6">System Overview</h3>
        
        {/* Overall status display */}
        <div className="text-center mb-6">
          <div className={`text-2xl font-bold mb-2 ${overallStatus.color}`}>
            {overallStatus.text}
          </div>
          <div className="text-white/70">
            {operationalCount}/{totalServices} services operational
          </div>
        </div>
        
        {/* Quick stats */}
        <div className="space-y-3">
          <div className="flex justify-between items-center p-2 rounded glass-secondary">
            <span className="text-white/70">Total Services</span>
            <span className="text-white font-medium">{totalServices}</span>
          </div>
          <div className="flex justify-between items-center p-2 rounded glass-secondary">
            <span className="text-white/70">Operational</span>
            <span className="text-green-400 font-medium">{operationalCount}</span>
          </div>
          {degradedCount > 0 && (
            <div className="flex justify-between items-center p-2 rounded glass-secondary">
              <span className="text-white/70">Degraded</span>
              <span className="text-yellow-400 font-medium">{degradedCount}</span>
            </div>
          )}
          {downCount > 0 && (
            <div className="flex justify-between items-center p-2 rounded glass-secondary">
              <span className="text-white/70">Down</span>
              <span className="text-red-400 font-medium">{downCount}</span>
            </div>
          )}
        </div>
      </div>
      
      {/* Service categories */}
      {categories.length > 0 && (
        <div className="glass-card rounded-2xl p-6 animate-desktop-fade-in">
          <h3 className="text-lg font-semibold text-white mb-4">Categories</h3>
          <div className="space-y-2">
            {categories.map(category => {
              const categoryStatus = getCategoryStatus(category);
              const categoryServices = services.filter(s => s.category === category);
              return (
                <div 
                  key={category} 
                  className="flex justify-between items-center p-2 rounded hover:bg-white/5 transition-colors cursor-pointer"
                >
                  <div>
                    <span className="text-white/80">{category}</span>
                    <div className="text-xs text-white/60">{categoryServices.length} services</div>
                  </div>
                  <span className={`text-lg ${categoryStatus.color}`}>
                    {categoryStatus.symbol}
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Quick actions */}
      <div className="glass-card rounded-2xl p-6 animate-desktop-fade-in">
        <h3 className="text-lg font-semibold text-white mb-4">Quick Actions</h3>
        <div className="space-y-2">
          <button className="w-full p-2 rounded glass-secondary hover:bg-white/10 transition-colors text-left text-white/80 text-sm">
            📊 View All Metrics
          </button>
          <button className="w-full p-2 rounded glass-secondary hover:bg-white/10 transition-colors text-left text-white/80 text-sm">
            🔄 Refresh All Services
          </button>
          <button className="w-full p-2 rounded glass-secondary hover:bg-white/10 transition-colors text-left text-white/80 text-sm">
            ⚙️ Service Configuration
          </button>
        </div>
      </div>
      
    </div>
  );
};

export default DesktopSidebar;