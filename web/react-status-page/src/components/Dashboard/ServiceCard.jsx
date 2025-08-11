import { useState } from 'react';
import StatusIndicator from './StatusIndicator';

const ServiceCard = ({ service }) => {
  const [isHovered, setIsHovered] = useState(false);

  // Handle missing or invalid service data
  if (!service || !service.name) {
    return (
      <div className="glass-card rounded-2xl p-6 min-h-[280px] lg:min-h-[320px] flex items-center justify-center">
        <div className="text-white/70 text-center">
          <div className="text-2xl mb-2">⚠️</div>
          <div>Invalid service data</div>
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

  const getStatusColors = (status) => {
    switch (status) {
      case 'operational':
        return {
          bg: 'bg-operational-glass',
          border: 'border-operational-border',
          text: 'text-green-300',
          dot: 'text-green-400'
        };
      case 'degraded':
        return {
          bg: 'bg-degraded-glass',
          border: 'border-degraded-border', 
          text: 'text-yellow-300',
          dot: 'text-yellow-400'
        };
      case 'down':
        return {
          bg: 'bg-down-glass',
          border: 'border-down-border',
          text: 'text-red-300',
          dot: 'text-red-400'
        };
      default:
        return {
          bg: 'bg-glass-secondary',
          border: 'border-glass-border',
          text: 'text-gray-300',
          dot: 'text-gray-400'
        };
    }
  };

  const statusColors = getStatusColors(service.status);

  return (
    <div
      className="relative glass-card rounded-2xl p-6 desktop-hover-lift cursor-pointer min-h-[280px] lg:min-h-[320px] animate-desktop-fade-in"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Enhanced header with larger elements for desktop */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-4">
          {/* Larger service icon for desktop */}
          <div className="w-12 h-12 lg:w-16 lg:h-16 rounded-full glass-card flex items-center justify-center text-2xl lg:text-3xl">
            {service.icon || '🌐'}
          </div>
          
          <div>
            {/* Larger service name for desktop */}
            <h3 className="text-xl lg:text-2xl font-semibold text-white">
              {service.name}
            </h3>
            {/* Service category for desktop */}
            <p className="text-sm text-white/70">{service.category || 'Service'}</p>
          </div>
        </div>
        
        {/* Enhanced status indicator with glassmorphism */}
        <div className={`
          px-4 py-2 rounded-full text-sm font-medium backdrop-blur-10 border
          ${statusColors.bg} ${statusColors.border} ${statusColors.text}
        `}>
          <span className={`inline-block w-2 h-2 rounded-full mr-2 ${statusColors.dot}`}>●</span>
          {service.status || 'operational'}
        </div>
      </div>
      
      {/* Desktop metrics grid */}
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="text-center p-3 rounded-lg glass-secondary backdrop-blur-10">
          <div className="text-2xl font-bold text-white">{service.uptime || '99.9'}%</div>
          <div className="text-xs text-white/70 uppercase tracking-wider">Uptime</div>
        </div>
        <div className="text-center p-3 rounded-lg glass-secondary backdrop-blur-10">
          <div className="text-2xl font-bold text-white">{formatLatency(service.latency_ms)}</div>
          <div className="text-xs text-white/70 uppercase tracking-wider">Response</div>
        </div>
      </div>
      
      {/* Enhanced details section */}
      <div className="space-y-3 text-sm">
        <div className="flex justify-between items-center p-2 rounded glass-secondary">
          <span className="text-white/70">Last Check:</span>
          <span className="text-white font-medium">
            {formatTimestamp(service.updated_at) || 'Just now'}
          </span>
        </div>
        
        {service.error && (
          <div className="p-3 bg-down-glass border border-down-border rounded-lg">
            <div className="text-red-300 text-xs font-medium mb-1">Error Details:</div>
            <div className="text-red-200 text-xs">{service.error}</div>
          </div>
        )}
      </div>
      
      {/* Desktop hover overlay */}
      {isHovered && (
        <div className="absolute inset-0 glass-card rounded-2xl p-6 flex items-center justify-center transition-opacity duration-200">
          <div className="text-center">
            <div className="text-3xl mb-3">📊</div>
            <p className="text-white font-medium mb-2">View Details</p>
            <p className="text-white/70 text-sm">Click to see full service metrics</p>
          </div>
        </div>
      )}
    </div>
  );
};

export default ServiceCard;