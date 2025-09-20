import { useState, useEffect } from 'react';
import usersApi from '../api/users';
import UserForm from './UserForm';
import LoadingSpinner from './LoadingSpinner';
import ErrorMessage from './ErrorMessage';

const UserList = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [editingUser, setEditingUser] = useState(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);

  useEffect(() => {
    loadUsers();
  }, [page]);

  const loadUsers = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await usersApi.getAll(page, pageSize);
      setUsers(data.users || []);
    } catch (err) {
      setError('Failed to load users');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (userData) => {
    try {
      await usersApi.create(userData);
      await loadUsers();
      setShowCreateForm(false);
    } catch (err) {
      setError('Failed to create user');
      console.error(err);
    }
  };

  const handleUpdate = async (id, userData) => {
    try {
      await usersApi.update(id, userData);
      await loadUsers();
      setEditingUser(null);
    } catch (err) {
      setError('Failed to update user');
      console.error(err);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this user?')) {
      return;
    }
    try {
      await usersApi.delete(id);
      await loadUsers();
    } catch (err) {
      setError('Failed to delete user');
      console.error(err);
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage message={error} />;

  return (
    <div className="glass-card rounded-lg overflow-hidden animate-fade-in">
      <div className="px-4 py-5 sm:px-6 flex justify-between items-center border-b border-rain-gray-border">
        <h3 className="text-lg font-medium text-white">Users</h3>
        <button
          onClick={() => setShowCreateForm(true)}
          className="btn-primary"
        >
          Add User
        </button>
      </div>

      {showCreateForm && (
        <div className="px-4 py-3 border-t border-rain-gray-border bg-rain-gray-medium/30">
          <UserForm
            onSubmit={handleCreate}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      )}

      <div>
        <ul className="divide-y divide-rain-gray-border">
          {users.map((user) => (
            <li key={user.id} className="px-4 py-4 sm:px-6 transition-all hover:bg-rain-gray-medium/20">
              {editingUser === user.id ? (
                <UserForm
                  user={user}
                  onSubmit={(data) => handleUpdate(user.id, data)}
                  onCancel={() => setEditingUser(null)}
                />
              ) : (
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-white">{user.name}</p>
                    <p className="text-sm text-gray-400">{user.email}</p>
                  </div>
                  <div className="flex space-x-3">
                    <button
                      onClick={() => setEditingUser(user.id)}
                      className="text-rain-pink hover:text-rain-pink-light text-sm transition-colors"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDelete(user.id)}
                      className="text-red-400 hover:text-red-300 text-sm transition-colors"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              )}
            </li>
          ))}
        </ul>
      </div>

      {users.length > 0 && (
        <div className="px-4 py-3 border-t border-rain-gray-border sm:px-6">
          <div className="flex justify-between items-center">
            <button
              onClick={() => setPage(Math.max(1, page - 1))}
              disabled={page === 1}
              className="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed text-sm"
            >
              Previous
            </button>
            <span className="text-sm text-gray-400">Page {page}</span>
            <button
              onClick={() => setPage(page + 1)}
              disabled={users.length < pageSize}
              className="btn-secondary disabled:opacity-50 disabled:cursor-not-allowed text-sm"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default UserList;