import React, { useState } from 'react';
import UserList from './components/UserList';
import UserForm from './components/UserForm';
import UserCard from './components/UserCard';
import { User } from './types/user';

type ViewMode = 'list' | 'create' | 'edit' | 'view';

function App() {
  const [currentView, setCurrentView] = useState<ViewMode>('list');
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleCreate = () => {
    setSelectedUser(null);
    setCurrentView('create');
  };

  const handleEdit = (user: User) => {
    setSelectedUser(user);
    setCurrentView('edit');
  };

  const handleView = (user: User) => {
    setSelectedUser(user);
    setCurrentView('view');
  };

  const handleFormSuccess = () => {
    setCurrentView('list');
    setSelectedUser(null);
    setRefreshTrigger(prev => prev + 1);
  };

  const handleCancel = () => {
    setCurrentView('list');
    setSelectedUser(null);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-rain-gray-50 via-white to-rain-purple-50">
      <header className="bg-white/80 backdrop-blur-md sticky top-0 z-40">
        {/* App Header */}
        <div className="border-b border-rain-gray-200">
          <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-xl bg-gradient-purple-pink flex items-center justify-center text-white shadow-lg">
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div>
                  <h1 className="text-xl font-bold text-gradient">Make It Rain</h1>
                  <p className="text-xs text-rain-gray-500">Rain Interview Platform</p>
                </div>
              </div>
              <div className="flex items-center gap-4">
                <span className="text-sm text-rain-gray-500">Welcome back</span>
                <div className="w-8 h-8 rounded-full bg-gradient-purple flex items-center justify-center text-white text-sm font-semibold">
                  A
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Navigation Tabs */}
        <div className="border-b border-rain-gray-200">
          <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
            <nav className="flex gap-8">
              <button className="relative py-3 text-sm font-medium text-rain-purple-600 border-b-2 border-rain-purple-600">
                <span className="flex items-center gap-2">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                  </svg>
                  Users
                </span>
              </button>
              {/* Easy to add new tabs during interview - just uncomment and modify:
              <button className="py-3 text-sm font-medium text-rain-gray-500 hover:text-rain-gray-700 transition-colors">
                <span className="flex items-center gap-2">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="..." />
                  </svg>
                  Tab Name
                </span>
              </button>
              */}
            </nav>
          </div>
        </div>

        {/* Page Header */}
        <div className="bg-white border-b border-rain-gray-100">
          <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
            <div className="flex justify-between items-center">
              <div>
                <h2 className="text-lg font-semibold text-rain-gray-900">User Management</h2>
                <p className="text-sm text-rain-gray-500 mt-0.5">Manage your platform users and permissions</p>
              </div>
              {currentView === 'list' && (
                <button
                  onClick={handleCreate}
                  className="group relative px-4 py-2 bg-gradient-purple-pink text-white rounded-lg font-medium shadow-lg shadow-rain-purple-500/20 hover:shadow-xl hover:shadow-rain-purple-500/30 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200 text-sm"
                >
                  <span className="flex items-center gap-2">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    Add User
                  </span>
                </button>
              )}
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <div className="animate-in">
          {currentView === 'list' && (
            <UserList
              onEdit={handleEdit}
              onView={handleView}
              refreshTrigger={refreshTrigger}
            />
          )}

          {(currentView === 'create' || currentView === 'edit') && (
            <div className="max-w-2xl mx-auto">
              <UserForm
                user={currentView === 'edit' ? selectedUser : null}
                onSuccess={handleFormSuccess}
                onCancel={handleCancel}
              />
            </div>
          )}

          {currentView === 'view' && selectedUser && (
            <div className="fixed inset-0 bg-black/20 backdrop-blur-sm flex items-center justify-center p-4 z-50 animate-in">
              <div className="max-w-2xl w-full">
                <UserCard
                  user={selectedUser}
                  onEdit={() => handleEdit(selectedUser)}
                  onClose={handleCancel}
                />
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

export default App;