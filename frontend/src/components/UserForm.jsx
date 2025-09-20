import { useState } from 'react';

const UserForm = ({ user, onSubmit, onCancel }) => {
  const [formData, setFormData] = useState({
    name: user?.name || '',
    email: user?.email || '',
    password: user ? '' : '',
  });
  const [errors, setErrors] = useState({});

  const validate = () => {
    const newErrors = {};
    if (!formData.name) newErrors.name = 'Name is required';
    if (!formData.email) newErrors.email = 'Email is required';
    else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'Email is invalid';
    }
    if (!user && !formData.password) {
      newErrors.password = 'Password is required';
    }
    return newErrors;
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const newErrors = validate();
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }
    const dataToSubmit = { ...formData };
    if (user && !formData.password) {
      delete dataToSubmit.password;
    }
    onSubmit(dataToSubmit);
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: '' }));
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-gray-300">
          Name
        </label>
        <input
          type="text"
          id="name"
          name="name"
          value={formData.name}
          onChange={handleChange}
          className="mt-1 block w-full rounded-md sm:text-sm px-3 py-2 bg-rain-gray-dark border border-rain-gray-border text-white placeholder-gray-500 focus:ring-2 focus:ring-rain-pink focus:ring-opacity-50 focus:border-transparent"
        />
        {errors.name && (
          <p className="mt-1 text-sm text-red-400">{errors.name}</p>
        )}
      </div>

      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-300">
          Email
        </label>
        <input
          type="email"
          id="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          className="mt-1 block w-full rounded-md sm:text-sm px-3 py-2 bg-rain-gray-dark border border-rain-gray-border text-white placeholder-gray-500 focus:ring-2 focus:ring-rain-pink focus:ring-opacity-50 focus:border-transparent"
        />
        {errors.email && (
          <p className="mt-1 text-sm text-red-400">{errors.email}</p>
        )}
      </div>

      {(!user || formData.password) && (
        <div>
          <label htmlFor="password" className="block text-sm font-medium text-gray-300">
            Password {user && <span className="text-gray-500">(leave blank to keep current)</span>}
          </label>
          <input
            type="password"
            id="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            className="mt-1 block w-full rounded-md sm:text-sm px-3 py-2 bg-rain-gray-dark border border-rain-gray-border text-white placeholder-gray-500 focus:ring-2 focus:ring-rain-pink focus:ring-opacity-50 focus:border-transparent"
          />
          {errors.password && (
            <p className="mt-1 text-sm text-red-400">{errors.password}</p>
          )}
        </div>
      )}

      <div className="flex justify-end space-x-3">
        <button
          type="button"
          onClick={onCancel}
          className="btn-secondary text-sm"
        >
          Cancel
        </button>
        <button
          type="submit"
          className="btn-primary text-sm"
        >
          {user ? 'Update' : 'Create'}
        </button>
      </div>
    </form>
  );
};

export default UserForm;