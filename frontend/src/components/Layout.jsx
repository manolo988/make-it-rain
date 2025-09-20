import { Link, Outlet, useLocation } from 'react-router-dom';

const Layout = () => {
  const location = useLocation();

  const isActive = (path) => {
    return location.pathname === path;
  };

  return (
    <div className="min-h-screen bg-gradient-dark">
      <nav className="glass-card border-b border-rain-gray-border">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <div className="flex-shrink-0 flex items-center">
                <h1 className="text-xl font-bold bg-gradient-rain bg-clip-text text-transparent animate-fade-in">
                  Make It Rain
                </h1>
              </div>
              <div className="hidden sm:ml-8 sm:flex sm:space-x-6">
                <Link
                  to="/auctions"
                  className={`inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-all duration-300 ${
                    location.pathname.startsWith('/auctions')
                      ? 'border-rain-pink text-rain-pink'
                      : 'border-transparent text-gray-400 hover:text-white hover:border-rain-pink/50'
                  }`}
                >
                  Auctions
                </Link>
                <Link
                  to="/users"
                  className={`inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium transition-all duration-300 ${
                    isActive('/users')
                      ? 'border-rain-pink text-rain-pink'
                      : 'border-transparent text-gray-400 hover:text-white hover:border-rain-pink/50'
                  }`}
                >
                  Users
                </Link>
              </div>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8 animate-fade-in">
        <Outlet />
      </main>
    </div>
  );
};

export default Layout;