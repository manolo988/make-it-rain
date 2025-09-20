import { useState, useEffect } from 'react';
import healthApi from '../api/health';

const HealthStatus = () => {
  const [status, setStatus] = useState({ health: null, ready: null });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkStatus();
    const interval = setInterval(checkStatus, 30000); // Check every 30 seconds
    return () => clearInterval(interval);
  }, []);

  const checkStatus = async () => {
    try {
      const [healthRes, readyRes] = await Promise.all([
        healthApi.check().catch(() => ({ status: 'unhealthy' })),
        healthApi.ready().catch(() => ({ status: 'not ready' })),
      ]);
      setStatus({
        health: healthRes.status === 'healthy',
        ready: readyRes.status === 'ready',
      });
    } catch (error) {
      console.error('Failed to check status:', error);
      setStatus({ health: false, ready: false });
    } finally {
      setLoading(false);
    }
  };

  const StatusIndicator = ({ label, isHealthy }) => (
    <div className="flex items-center space-x-2">
      <div
        className={`h-3 w-3 rounded-full ${
          isHealthy ? 'bg-green-500 shadow-green-500/50' : 'bg-red-500 shadow-red-500/50'
        } ${!loading && 'animate-pulse'} shadow-lg`}
      />
      <span className="text-sm text-gray-400">{label}</span>
    </div>
  );

  if (loading) {
    return (
      <div className="glass-card rounded-lg p-4 h-full">
        <div className="flex items-center space-x-4">
          <div className="h-3 w-3 rounded-full bg-gray-600 animate-pulse" />
          <span className="text-sm text-gray-400">Checking status...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="glass-card rounded-lg p-4 h-full card-hover">
      <div className="relative">
        <div className="absolute inset-0 bg-gradient-to-br from-rain-pink/10 to-transparent rounded-lg"></div>
        <div className="relative">
          <h3 className="text-sm font-medium text-white mb-3">API Status</h3>
          <div className="space-y-2">
            <StatusIndicator label="Health Check" isHealthy={status.health} />
            <StatusIndicator label="Readiness" isHealthy={status.ready} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default HealthStatus;