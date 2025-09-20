import { useState, useEffect } from 'react';
import HealthStatus from '../components/HealthStatus';
import usersApi from '../api/users';

const Dashboard = () => {
  const [stats, setStats] = useState({
    totalUsers: 0,
    recentUsers: [],
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      const usersData = await usersApi.getAll(1, 5);
      setStats({
        totalUsers: usersData.total || usersData.users?.length || 0,
        recentUsers: usersData.users || [],
      });
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="px-4 sm:px-0">
      <h1 className="text-3xl font-bold mb-6 bg-gradient-rain bg-clip-text text-transparent animate-slide-up">
        Dashboard
      </h1>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div className="glass-card rounded-lg overflow-hidden card-hover animate-slide-up">
          <div className="px-4 py-5 sm:p-6 relative">
            <div className="absolute inset-0 bg-gradient-to-br from-rain-pink/10 to-transparent"></div>
            <div className="relative">
              <dt className="text-sm font-medium text-gray-400 truncate">
                Total Users
              </dt>
              <dd className="mt-1 text-3xl font-bold text-white">
                {loading ? '...' : stats.totalUsers}
              </dd>
            </div>
          </div>
        </div>

        <div className="glass-card rounded-lg overflow-hidden card-hover animate-slide-up" style={{animationDelay: '0.1s'}}>
          <div className="px-4 py-5 sm:p-6 relative">
            <div className="absolute inset-0 bg-gradient-to-br from-rain-pink/10 to-transparent"></div>
            <div className="relative">
              <dt className="text-sm font-medium text-gray-400 truncate">
                API Version
              </dt>
              <dd className="mt-1 text-3xl font-bold text-white">
                v1
              </dd>
            </div>
          </div>
        </div>

        <div className="sm:col-span-2 lg:col-span-1 animate-slide-up" style={{animationDelay: '0.2s'}}>
          <HealthStatus />
        </div>
      </div>

      <div className="mt-8 animate-slide-up" style={{animationDelay: '0.3s'}}>
        <div className="glass-card rounded-lg overflow-hidden">
          <div className="px-4 py-5 sm:px-6 border-b border-rain-gray-border">
            <h3 className="text-lg font-medium text-white">
              Recent Users
            </h3>
          </div>
          <div>
            {loading ? (
              <div className="px-4 py-4 text-sm text-gray-400">Loading...</div>
            ) : stats.recentUsers.length > 0 ? (
              <ul className="divide-y divide-rain-gray-border">
                {stats.recentUsers.map((user) => (
                  <li key={user.id} className="px-4 py-4 sm:px-6 transition-colors hover:bg-rain-gray-medium/30">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-sm font-medium text-white">
                          {user.name}
                        </p>
                        <p className="text-sm text-gray-400">{user.email}</p>
                      </div>
                      <div className="text-sm text-rain-pink/70">
                        ID: {user.id}
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            ) : (
              <div className="px-4 py-4 text-sm text-gray-400">
                No users found
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;