import React from 'react';
import { User } from '../types/user';

interface UserCardProps {
  user: User;
  onEdit: () => void;
  onClose: () => void;
}

const UserCard: React.FC<UserCardProps> = ({ user, onEdit, onClose }) => {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className="card p-8 shadow-soft relative max-h-[90vh] overflow-y-auto">
      <button
        onClick={onClose}
        className="absolute top-6 right-6 p-2 text-rain-gray-400 hover:text-rain-gray-600 hover:bg-rain-gray-100 rounded-lg transition-all"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <div className="mb-8">
        <div className="flex items-center gap-4 mb-6">
          <div className="w-20 h-20 rounded-2xl bg-gradient-purple-pink flex items-center justify-center text-white font-bold text-3xl shadow-lg">
            {user.name.charAt(0).toUpperCase()}
          </div>
          <div>
            <h2 className="text-2xl font-bold text-rain-gray-900">{user.name}</h2>
            <p className="text-rain-gray-500">{user.email}</p>
          </div>
        </div>
      </div>

      <div className="space-y-6">
        <div className="grid grid-cols-2 gap-6">
          <div className="space-y-1">
            <label className="text-xs font-semibold text-rain-gray-500 uppercase tracking-wider">User ID</label>
            <p className="text-lg font-medium text-rain-gray-900">#{user.id}</p>
          </div>

          <div className="space-y-1">
            <label className="text-xs font-semibold text-rain-gray-500 uppercase tracking-wider">Status</label>
            <div>
              <span className={user.is_active ? 'badge-active' : 'badge-inactive'}>
                {user.is_active ? 'Active' : 'Inactive'}
              </span>
            </div>
          </div>
        </div>

        <div className="pt-6 border-t border-rain-gray-100">
          <div className="grid grid-cols-1 gap-4">
            <div className="space-y-1">
              <label className="text-xs font-semibold text-rain-gray-500 uppercase tracking-wider">Email Address</label>
              <p className="text-base text-rain-gray-900 font-medium">{user.email}</p>
            </div>

            <div className="space-y-1">
              <label className="text-xs font-semibold text-rain-gray-500 uppercase tracking-wider">Full Name</label>
              <p className="text-base text-rain-gray-900 font-medium">{user.name}</p>
            </div>
          </div>
        </div>

        <div className="pt-6 border-t border-rain-gray-100">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-1">
              <label className="text-xs font-semibold text-rain-gray-500 uppercase tracking-wider">Created</label>
              <p className="text-sm text-rain-gray-700">{formatDate(user.created_at)}</p>
            </div>

            <div className="space-y-1">
              <label className="text-xs font-semibold text-rain-gray-500 uppercase tracking-wider">Last Updated</label>
              <p className="text-sm text-rain-gray-700">{formatDate(user.updated_at)}</p>
            </div>
          </div>
        </div>
      </div>

      <div className="flex gap-3 mt-8 pt-6 border-t border-rain-gray-100">
        <button
          onClick={onEdit}
          className="flex-1 btn-primary"
        >
          <span className="flex items-center justify-center gap-2">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            Edit User
          </span>
        </button>
        <button
          onClick={onClose}
          className="flex-1 btn-secondary"
        >
          Close
        </button>
      </div>
    </div>
  );
};

export default UserCard;